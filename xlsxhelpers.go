package main

import "github.com/tealeg/xlsx/v3"

var (
	style1, style3, style3Left, style4, style4Right, style5, style6, style6WithoutBold, styleBlueColor, styleBlueColorLeft, styleRedColor, styleRedColorLeft, styleHyperLink, styleHyperLinkAv0, styleHyperLinkContendersIsJury, styleHyperLinkBlueColor *xlsx.Style
)

func init() {
	defaultBorder := xlsx.Border{Left: "thin", LeftColor: "AAAAAA", Right: "thin", RightColor: "AAAAAA", Bottom: "thin", BottomColor: "AAAAAA", Top: "thin", TopColor: "AAAAAA"}

	style1 = xlsx.NewStyle()
	style1.Font.Name = "Arial"
	style1.Font.Size = 24
	style1.Border = defaultBorder

	style3 = xlsx.NewStyle()
	style3.Font.Name = "Arial"
	style3.Font.Size = 10
	style3.Font.Color = "FFFFFF"
	style3.Font.Bold = true
	style3.Fill.FgColor = "5B95F9"
	style3.Fill.PatternType = "solid"
	style3.Alignment.Horizontal = "center"
	style3.Border = defaultBorder

	style3Left = xlsx.NewStyle()
	style3Left.Font.Name = "Arial"
	style3Left.Font.Size = 10
	style3Left.Font.Color = "FFFFFF"
	style3Left.Font.Bold = true
	style3Left.Fill.FgColor = "5B95F9"
	style3Left.Fill.PatternType = "solid"
	style3Left.Alignment.Horizontal = "left"
	style3Left.Border = defaultBorder

	style4 = xlsx.NewStyle()
	style4.Font.Name = "Arial"
	style4.Font.Size = 10
	style4.Font.Bold = true
	style4.Fill.FgColor = "ACC9FE"
	style4.Alignment.Horizontal = "center"
	style4.Fill.PatternType = "solid"
	style4.Border = defaultBorder

	style4Right = xlsx.NewStyle()
	style4Right.Font.Name = "Arial"
	style4Right.Font.Size = 10
	style4Right.Font.Bold = true
	style4Right.Fill.FgColor = "ACC9FE"
	style4Right.Alignment.Horizontal = "right"
	style4Right.Fill.PatternType = "solid"
	style4Right.Border = defaultBorder

	style5 = xlsx.NewStyle()
	style5.Font.Name = "Arial"
	style5.Font.Size = 10
	style5.Font.Color = "FFFFFF"
	style5.Font.Bold = true
	style5.Fill.FgColor = "5B95F9"
	style5.Fill.PatternType = "solid"
	style5.Alignment.Horizontal = "center"
	style5.Border = defaultBorder

	style6 = xlsx.NewStyle()
	style6.Font.Name = "Arial"
	style6.Font.Size = 10
	style6.Font.Bold = true
	style6.Alignment.Horizontal = "center"
	style6.Border = defaultBorder

	style6WithoutBold = xlsx.NewStyle()
	style6WithoutBold.Font.Name = "Arial"
	style6WithoutBold.Font.Size = 10
	style6WithoutBold.Alignment.Horizontal = "left"
	style6WithoutBold.Border = defaultBorder

	styleBlueColor = xlsx.NewStyle()
	styleBlueColor.Font.Name = "Arial"
	styleBlueColor.Font.Size = 10
	styleBlueColor.Font.Bold = true
	styleBlueColor.Alignment.Horizontal = "center"
	styleBlueColor.Fill.FgColor = "E8F0FE"
	styleBlueColor.Fill.PatternType = "solid"
	styleBlueColor.Border = defaultBorder

	styleBlueColorLeft = xlsx.NewStyle()
	styleBlueColorLeft.Font.Name = "Arial"
	styleBlueColorLeft.Font.Size = 10
	styleBlueColorLeft.Alignment.Horizontal = "left"
	styleBlueColorLeft.Fill.FgColor = "E8F0FE"
	styleBlueColorLeft.Fill.PatternType = "solid"
	styleBlueColorLeft.Border = defaultBorder

	styleRedColor = xlsx.NewStyle()
	styleRedColor.Font.Name = "Arial"
	styleRedColor.Font.Size = 10
	styleRedColor.Font.Bold = true
	styleRedColor.Alignment.Horizontal = "center"
	styleRedColor.Fill.FgColor = "F4CCCC"
	styleRedColor.Fill.PatternType = "solid"
	styleRedColor.Border = defaultBorder

	styleRedColorLeft = xlsx.NewStyle()
	styleRedColorLeft.Font.Name = "Arial"
	styleRedColorLeft.Font.Size = 10
	styleRedColorLeft.Alignment.Horizontal = "left"
	styleRedColorLeft.Fill.FgColor = "F4CCCC"
	styleRedColorLeft.Fill.PatternType = "solid"
	styleRedColorLeft.Border = defaultBorder

	styleHyperLink = xlsx.NewStyle()
	styleHyperLink.Font.Name = "Arial"
	styleHyperLink.Font.Size = 10
	styleHyperLink.Font.Color = "1155CC"
	styleHyperLink.Font.Underline = true
	styleHyperLink.Border = defaultBorder

	styleHyperLinkAv0 = xlsx.NewStyle()
	styleHyperLinkAv0.Font.Name = "Arial"
	styleHyperLinkAv0.Font.Size = 10
	styleHyperLinkAv0.Font.Color = "1155CC"
	styleHyperLinkAv0.Font.Underline = true
	styleHyperLinkAv0.Border = defaultBorder
	styleHyperLinkAv0.Fill.FgColor = "F4CCCC"
	styleHyperLinkAv0.Fill.PatternType = "solid"

	styleHyperLinkContendersIsJury = xlsx.NewStyle()
	styleHyperLinkContendersIsJury.Font.Name = "Arial"
	styleHyperLinkContendersIsJury.Font.Size = 10
	styleHyperLinkContendersIsJury.Font.Color = "1155CC"
	styleHyperLinkContendersIsJury.Font.Underline = true
	styleHyperLinkContendersIsJury.Border = defaultBorder
	styleHyperLinkContendersIsJury.Fill.FgColor = "00FF00"
	styleHyperLinkContendersIsJury.Fill.PatternType = "solid"

	styleHyperLinkBlueColor = xlsx.NewStyle()
	styleHyperLinkBlueColor.Font.Name = "Arial"
	styleHyperLinkBlueColor.Font.Size = 10
	styleHyperLinkBlueColor.Font.Color = "1155CC"
	styleHyperLinkBlueColor.Font.Underline = true
	styleHyperLinkBlueColor.Border = defaultBorder
	styleHyperLinkBlueColor.Fill.FgColor = "E8F0FE"
	styleHyperLinkBlueColor.Fill.PatternType = "solid"
}

func addEmptyString(sheet *xlsx.Sheet, row, col int) {
	rowN := sheet.AddRow()
	cell1R3, _ := rowN.Sheet.Cell(row, col)
	cell1R3.String()
}

func addEmptyCellWithoutStyle(row *xlsx.Row) *xlsx.Cell {
	cell := row.AddCell()
	cell.SetString(" ")
	return cell
}

func addEmptyCellWithStyle(row *xlsx.Row, style *xlsx.Style) *xlsx.Cell {
	cell := row.AddCell()
	cell.SetStyle(style)
	cell.SetString(" ")
	return cell
}

func addCell(row *xlsx.Row, style *xlsx.Style, valueField string) *xlsx.Cell {
	cell := row.AddCell()
	cell.SetStyle(style)
	cell.SetString(valueField)
	return cell
}
