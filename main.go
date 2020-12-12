package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	goton "github.com/move-ton/ton-client-go"
	"github.com/move-ton/ton-client-go/domain"
	"github.com/tealeg/xlsx/v3"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Addresses not entered!")
	}

	ton, err := goton.NewTon(2)
	if err != nil {
		log.Fatal(err)
	}

	defer ton.Client.Destroy()

	datesAddress := os.Args[1:]

	file, err := os.Open("FreeTonContest.abi.json")
	if err != nil {
		fmt.Println("Error 1 open: ", err)
		return
	}

	abiByte, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error 2 read file abi: ", err)
		return
	}

	abiContract := domain.AbiContract{}
	err = json.Unmarshal(abiByte, &abiContract)
	if err != nil {
		fmt.Println(err)
		return
	}

	abiParams := domain.NewAbiSerialized()
	abiParams.Value = abiContract

	wb := xlsx.NewFile()

	for _, address := range datesAddress {
		fmt.Println("Get data to address: ", address)
		result, err := ton.Net.WaitForCollection(domain.ParamsOfWaitForCollection{
			Collection: "accounts",
			Filter:     json.RawMessage(fmt.Sprintf(`{"id":{"eq":"%s"}}`, address)),
			Result:     "id, boc"})
		if err != nil {
			fmt.Println("Bad query for boc data: ", err)
			continue
		}
		var (
			objmap map[string]json.RawMessage
			boc    string
		)

		err = json.Unmarshal(result.Result, &objmap)
		if err != nil {
			fmt.Println("Error unmarshalling result from graphql: ", err)
			continue
		}
		err = json.Unmarshal(objmap["boc"], &boc)
		if err != nil {
			fmt.Println("Error unmarshalling boc field from result: ", err)
			continue
		}

		encodingParams := domain.NewMessageSourceEncodingParams()
		encodingParams.Abi = abiParams
		encodingParams.Address = address
		encodingParams.CallSet = &domain.CallSet{FunctionName: "listContenders"}
		encodingParams.Signer = domain.NewSignerNone()
		message, err := ton.Abi.EncodeMessage(encodingParams)
		if err != nil {
			fmt.Println("Error Encode Message: ", err)
			continue
		}

		resultReq, err := ton.Tvm.RunTvm(domain.ParamsOfRunTvm{Message: message.Message, Account: boc, Abi: abiParams})
		if err != nil {
			fmt.Println("Error RunTvm for listContenders: ", err)
			continue
		}

		resultSubmission := &resultContenders{}
		err = json.Unmarshal(resultReq.Decoded.Output, resultSubmission)
		if err != nil {
			fmt.Println("Error unmarshalling result contenders: ", err)
			continue
		}

		md := &mainDats{}
		listContenders := make(map[string]int64)
		lenC := len(resultSubmission.Addresses)
		for i := 0; i < lenC; i++ {
			contDrs := contenders{}
			contDrs.IDS, _ = strconv.ParseInt(resultSubmission.Ids[i], 0, 64)
			contDrs.Address = resultSubmission.Addresses[i]

			md.Contenders = append(md.Contenders, contDrs)
			listContenders[contDrs.Address] = contDrs.IDS
		}

		encodingParams.CallSet.FunctionName = "getContestInfo"
		message, err = ton.Abi.EncodeMessage(encodingParams)
		if err != nil {
			fmt.Println("Error Encode Message: ", err)
			continue
		}

		resultReq, err = ton.Tvm.RunTvm(domain.ParamsOfRunTvm{Message: message.Message, Account: boc, Abi: abiParams})
		if err != nil {
			fmt.Println("Error RunTvm for getContestInfo: ", err)
			continue
		}

		res2 := &resContestInfo{}
		err = json.Unmarshal(resultReq.Decoded.Output, res2)
		if err != nil {
			fmt.Println("Error unmarshalling contest info: ", err)
			continue
		}

		md.TitleContext = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(hexToString([]byte(res2.Title)), "Contest Proposal: ", ""), "Contest: ", ""))
		md.LinkToContext = hexToString([]byte(res2.Link))
		lenJ := len(res2.JuryKeys)
		var slJury []string
		for i := 0; i < lenJ; i++ {
			juryS := jury{}
			juryS.Address = res2.JuryAddresses[i]
			juryS.PublicKey = res2.JuryKeys[i]

			md.Jurys = append(md.Jurys, juryS)
			slJury = append(slJury, juryS.Address)
		}

		mm := make(map[string]votes)
		votesTable := make(map[int64]map[string]allVotes)
		encodingParams.CallSet.FunctionName = "getVotesPerJuror"
		for n, val := range md.Contenders {
			if sort.SearchStrings(slJury, val.Address) > 0 {
				val.Jury = true
			}

			idReq := req{ID: val.IDS}
			encodingParams.CallSet.Input = idReq
			message, err = ton.Abi.EncodeMessage(encodingParams)
			if err != nil {
				fmt.Println("Error Encode Message: ", err)
				continue
			}
			resultReq, err = ton.Tvm.RunTvm(domain.ParamsOfRunTvm{Message: message.Message, Account: boc, Abi: abiParams})
			if err != nil {
				fmt.Println("Error RunTvm for getVotesPerJuror: ", err)
				continue
			}

			res3 := &goverment{}
			_ = json.Unmarshal(resultReq.Decoded.Output, res3)

			md.Contenders[n].GovermentD = res3

			var totalFF, sumCount int64
			for indx, valFor := range res3.JurorsFor {
				if vvF, found := mm[valFor]; found {
					vvF.JuryFor++
					mm[valFor] = vvF
				} else {
					mm[valFor] = votes{1, 0, 0}
				}

				ff, err := strconv.ParseInt(res3.Marks[indx], 0, 16)
				if err != nil {
					continue
				}
				if valNow, ok := votesTable[val.IDS]; ok {
					valNow[valFor] = allVotes{Mark: ff}
				} else {
					votesTable[val.IDS] = map[string]allVotes{
						valFor: allVotes{Mark: ff},
					}
				}

				totalFF += ff
				sumCount++
			}

			for _, valAbs := range res3.JurorsAbstained {
				if vvAb, found := mm[valAbs]; found {
					vvAb.JuryAbstained++
					mm[valAbs] = vvAb
				} else {
					mm[valAbs] = votes{0, 1, 0}
				}

				if valNow, ok := votesTable[val.IDS]; ok {
					valNow[valAbs] = allVotes{Mark: 0, Abstain: true}
				} else {
					votesTable[val.IDS] = map[string]allVotes{
						valAbs: allVotes{Mark: 0, Abstain: true},
					}
				}
			}

			for _, valAg := range res3.JurorsAgainst {
				if vvAg, found := mm[valAg]; found {
					vvAg.JuryAgainst++
					mm[valAg] = vvAg
				} else {
					mm[valAg] = votes{0, 0, 1}
				}

				if valNow, ok := votesTable[val.IDS]; ok {
					valNow[valAg] = allVotes{Mark: 0}
				} else {
					votesTable[val.IDS] = map[string]allVotes{
						valAg: allVotes{Mark: 0},
					}
				}
			}

			avScore := float64(totalFF) / float64(sumCount)
			if math.IsNaN(avScore) {
				md.Contenders[n].AverageScore = 0
			} else {
				md.Contenders[n].AverageScore = avScore
			}

			md.Contenders[n].Reject = int64(len(res3.JurorsAgainst))
		}
		generateFile(wb, md, mm, listContenders, votesTable)
	}

	generateResult(wb)
	err = wb.Save("ResultContests.xlsx")
	if err != nil {
		log.Fatal("Error save file: " + err.Error())
	}

	fmt.Println("Create file: ", "ResultContests.xlsx")

}

func generateFile(wb *xlsx.File, data *mainDats, mm map[string]votes, listContenders map[string]int64, votesTable map[int64]map[string]allVotes) {
	fmt.Println("Add to file contest: ", data.TitleContext)
	nameFile := data.TitleContext
	if len(nameFile) > 29 {
		nameFile = data.TitleContext[:28]
	}
	nameFile = strings.ReplaceAll(nameFile, "\\", " ")
	nameFile = strings.ReplaceAll(nameFile, "/", " ")
	nameFile = strings.ReplaceAll(nameFile, "?", " ")
	nameFile = strings.ReplaceAll(nameFile, "*", " ")
	nameFile = strings.ReplaceAll(nameFile, "[", " ")
	nameFile = strings.ReplaceAll(nameFile, "]", " ")
	nameFile = strings.ReplaceAll(nameFile, ":", " ")

	ch := 1
	nameFileTek := nameFile
	for {
		if _, ok := wb.Sheet[nameFileTek]; !ok {
			nameFile = nameFileTek
			break
		} else {
			nameFileTek = nameFile + strconv.Itoa(ch)
		}

		ch++
	}
	sheet1, err := wb.AddSheet(nameFile)
	if err != nil {
		fmt.Println(nameFile)
		log.Fatal(err)
	}
	sheet1.SetColWidth(0, 0, 70)
	sheet1.SetColWidth(1, 100, 15)

	addEmptyString(sheet1, 0, 0)

	row2 := sheet1.AddRow()
	cell1R2 := row2.AddCell()
	cell1R2.SetHyperlink(data.LinkToContext, data.TitleContext, "")
	cell1R2.SetStyle(style1)
	cell1R2.GetStyle().Font.Color = "1155CC"
	cell1R2.GetStyle().Font.Bold = true
	cell1R2.GetStyle().Font.Underline = true

	addEmptyCellWithoutStyle(row2)
	addEmptyCellWithoutStyle(row2)
	addEmptyCellWithoutStyle(row2)

	addCell(row2, style1, "Jurors Votes").GetStyle().Font.Underline = false

	addEmptyString(sheet1, 2, 0)

	row4 := sheet1.AddRow()

	addCell(row4, style3, "Wallet Address")
	addCell(row4, style3, "Submission №")
	addCell(row4, style3, "Average score")

	addEmptyCellWithStyle(row4, style3)

	for _, val := range data.Jurys {
		addCell(row4, style3Left, val.Address)
		addCell(row4, style3, "Score Diff")
	}

	var (
		idx       int64
		blueColor bool
	)

	idx = 1
	jurysAverageDiff := make(map[string]float64)
	for _, val := range data.Contenders {
		row5 := sheet1.AddRow()
		cell1R5 := addEmptyCellWithStyle(row5, styleHyperLink)
		cell1R5.SetHyperlink(linkToExplorer+val.Address, val.Address, "")
		if val.AverageScore == 0 {
			cell1R5.SetStyle(styleHyperLinkAv0)
		} else if val.Jury {
			cell1R5.SetStyle(styleHyperLinkContendersIsJury)
		} else if blueColor {
			cell1R5.SetStyle(styleHyperLinkBlueColor)
		}

		cell2R5 := row5.AddCell()
		cell2R5.SetValue(val.IDS)
		cell2R5.SetStyle(style6)
		if val.AverageScore == 0 {
			cell2R5.SetStyle(styleRedColor)
		} else if blueColor {
			cell2R5.SetStyle(styleBlueColor)
		}

		cell3R5 := row5.AddCell()
		cell3R5.SetFloatWithFormat(val.AverageScore, "#0.00")
		cell3R5.SetStyle(style6)
		if val.AverageScore == 0 {
			cell3R5.SetStyle(styleRedColor)
		} else if blueColor {
			cell3R5.SetStyle(styleBlueColor)
		}

		if !blueColor {
			addEmptyCellWithStyle(row5, style6)
		} else {
			addEmptyCellWithStyle(row5, styleBlueColor)
		}

		for _, val2 := range data.Jurys {
			cell5R5 := addEmptyCellWithStyle(row5, style6)
			cell6R5 := addEmptyCellWithStyle(row5, style6)
			if valN, ok := votesTable[val.IDS][val2.Address]; ok {
				if !valN.Abstain {
					cell5R5.SetInt64(valN.Mark)

					valSD := float64(valN.Mark) - val.AverageScore
					if valSD < 0 {
						valSD = valSD * (-1)
					}
					cell6R5.SetFloatWithFormat(valSD, "#0.00")

					if _, ok := jurysAverageDiff[val2.Address]; ok {
						jurysAverageDiff[val2.Address] += valSD
					} else {
						jurysAverageDiff[val2.Address] = valSD
					}
				} else {
					cell5R5.SetString("Abstained")
					cell6R5.SetFloatWithFormat(0, "#0.00")
				}
			} else {
				cell5R5.SetString("No Vote")
				cell6R5.SetFloatWithFormat(0, "#0.00")
			}

			if blueColor {
				cell5R5.SetStyle(styleBlueColor)
				cell6R5.SetStyle(styleBlueColor)
			}
		}

		blueColor = !blueColor
		idx++
	}

	row6 := sheet1.AddRow()
	addCell(row6, style4Right, "Total:")
	addEmptyCellWithStyle(row6, style4).SetInt(len(data.Contenders))
	addEmptyCellWithStyle(row6, style4)
	addEmptyCellWithStyle(row6, style4)

	countJur := len(data.Jurys)
	avgDiffJurys := make(map[string]float64)
	for i := 0; i < countJur; i++ {
		cell5R6 := row6.AddCell()
		cell5R6.SetStyle(style4)
		cell5R6.SetString("Avg. Diff")

		cell6R6 := row6.AddCell()
		cell6R6.SetStyle(style4)
		datAvgDiff := jurysAverageDiff[data.Jurys[i].Address]
		datmm := mm[data.Jurys[i].Address]
		avgDiff := datAvgDiff / float64(datmm.JuryFor+datmm.JuryAgainst+datmm.JuryAbstained)
		if avgDiff != 0 && !math.IsNaN(avgDiff) {
			cell6R6.SetFloatWithFormat(avgDiff, "#0.00")
		} else {
			cell6R6.SetFloatWithFormat(0.00, "#0.00")
		}
		avgDiffJurys[data.Jurys[i].Address] = avgDiff
	}

	addEmptyString(sheet1, 2, 0)

	row7 := sheet1.AddRow()
	addCell(row7, style1, "Jury Activity").GetStyle().Font.Bold = true

	addEmptyString(sheet1, 2, 0)

	row8 := sheet1.AddRow()
	addCell(row8, style5, "Wallet Address")
	addCell(row8, style5, "Jury №")
	addCell(row8, style5, "Votes count")
	addCell(row8, style5, "Abstained")
	addCell(row8, style5, "Participant?")
	addCell(row8, style5, "Efficiency, %")
	addCell(row8, style5, "Avg. Diff")

	blueColor = false
	indJury := 1
	for _, valJ := range data.Jurys {
		row9 := sheet1.AddRow()
		cell1R9 := addCell(row9, style6WithoutBold, valJ.Address)
		if blueColor {
			cell1R9.SetStyle(styleBlueColorLeft)
		}

		cell2R9 := addEmptyCellWithStyle(row9, style6)
		cell2R9.SetValue(indJury)
		if blueColor {
			cell2R9.SetStyle(styleBlueColor)
		}

		dateVotes := mm[valJ.Address]
		cell3R9 := addEmptyCellWithStyle(row9, style6)
		cell3R9.SetInt64(dateVotes.JuryAbstained + dateVotes.JuryAgainst + dateVotes.JuryFor)
		if blueColor {
			cell3R9.SetStyle(styleBlueColor)
		}

		cell4R9 := addEmptyCellWithStyle(row9, style6)
		cell4R9.SetInt64(dateVotes.JuryAbstained)
		if blueColor {
			cell4R9.SetStyle(styleBlueColor)
		}

		cell5R9 := addEmptyCellWithStyle(row9, style6)
		var juryContenders bool
		if value, found := listContenders[valJ.Address]; found {
			cell5R9.SetInt64(value)
			juryContenders = true
		} else {
			cell5R9.SetString("No")
		}

		if blueColor {
			cell5R9.SetStyle(styleBlueColor)
		}

		efficiency := 0.0
		if juryContenders {
			efficiency = 100
		} else {
			if (float64(dateVotes.JuryAbstained) / float64(len(data.Contenders))) < 0.2 {
				efficiency = ((float64(dateVotes.JuryAbstained+dateVotes.JuryAgainst+dateVotes.JuryFor) / float64(len(data.Contenders))) * 100)
			} else {
				efficiency = ((float64(dateVotes.JuryAgainst+dateVotes.JuryFor) / float64(len(data.Contenders))) * 100)
			}
		}

		cell6R9 := addEmptyCellWithStyle(row9, style6)
		cell6R9.SetFloatWithFormat(efficiency, "#0.00")
		if blueColor {
			cell6R9.SetStyle(styleBlueColor)
		}

		valAvgM := 0.00
		cell7R9 := row9.AddCell()
		if valAvg, ok := avgDiffJurys[valJ.Address]; ok && !math.IsNaN(valAvg) && valAvg != 0 {
			valAvgM = valAvg
		}
		cell7R9.SetFloatWithFormat(valAvgM, "#0.00")
		cell7R9.SetStyle(style6)
		if blueColor {
			cell7R9.SetStyle(styleBlueColor)
		}

		if efficiency == 0 {
			cell1R9.SetStyle(styleRedColorLeft)
			cell2R9.SetStyle(styleRedColor)
			cell3R9.SetStyle(styleRedColor)
			cell4R9.SetStyle(styleRedColor)
			cell5R9.SetStyle(styleRedColor)
			cell6R9.SetStyle(styleRedColor)
			cell7R9.SetStyle(styleRedColor)
		}

		if valNow, ok := resultData[valJ.Address]; ok {
			valNow.Efficiency += efficiency
			valNow.CountEfficiency++
			valNow.AvgDiff += valAvgM

			if dateVotes.JuryAbstained+dateVotes.JuryAgainst+dateVotes.JuryFor != 0 {
				valNow.CountAvgDiff++
			}

			resultData[valJ.Address] = valNow
		} else {
			resultData[valJ.Address] = statJurys{
				Efficiency:      efficiency,
				CountEfficiency: 1,
				AvgDiff:         valAvgM,
				CountAvgDiff:    1,
			}
		}

		indJury++
		blueColor = !blueColor
	}
}

func generateResult(wb *xlsx.File) {
	sheet1, err := wb.AddSheet("TOTAL")
	if err != nil {
		log.Fatal(err)
	}
	sheet1.SetColWidth(0, 0, 70)
	sheet1.SetColWidth(1, 10, 15)

	addEmptyString(sheet1, 0, 0)
	addCell(sheet1.AddRow(), style1, "Overall Jury Efficiency")
	addEmptyString(sheet1, 0, 0)

	row1 := sheet1.AddRow()
	addCell(row1, style3, "Wallet Address")
	addCell(row1, style3, "Efficiency, %")
	addCell(row1, style3, "Avg. Diff")

	resultStore := make(map[string]resultTable)

	for address, value := range resultData {
		resultStore[address] = resultTable{
			address:    address,
			efficiency: value.Efficiency / value.CountEfficiency,
			avgDiff:    value.AvgDiff / value.CountAvgDiff,
		}
	}

	s := make(dataSlice, 0, len(resultStore))

	for _, d := range resultStore {
		s = append(s, d)
	}

	sort.Sort(s)

	var blueColor bool

	for _, value := range s {
		row2 := sheet1.AddRow()
		cell1R2 := addCell(row2, style6WithoutBold, value.address)
		if blueColor {
			cell1R2.SetStyle(styleBlueColorLeft)
		}

		cell2R2 := addEmptyCellWithStyle(row2, style6)
		cell2R2.SetFloatWithFormat(value.efficiency, "#0.00")

		cell3R2 := addEmptyCellWithStyle(row2, style6)
		cell3R2.SetFloatWithFormat(value.avgDiff, "#0.00")
		if blueColor {
			cell3R2.SetStyle(styleBlueColor)
			cell2R2.SetStyle(styleBlueColor)
		}
		blueColor = !blueColor
	}
}
