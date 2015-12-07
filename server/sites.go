package server

import (
	"github.com/siim-/siil/entity/site"
)

func dataMarkup(dt, dd string) string {
	var str string = ""
	str += surroundWithDt(dt)
	str += surroundWithDd(dd)
	return str
}

func surroundWithH3(str string) string {
	return "<h3>" + str + "</h3>"
}

func surroundWithDl(str string) string {
	return "<dl>" + str + "</dl>"
}

func surroundWithDt(str string) string {
	return "<dt>" + str + "</dt>"
}

func surroundWithDd(str string) string {
	return "<dd>" + str + "</dd>"
}

func surroundWithSection(str string) string {
	return "<section>" + str + "</section>"
}

func surroundWithColumn(str, colWidth string) string {
	return "<div class=\"" + colWidth + " columns\">" + str + "</div>"
}

func surroundWithRow(str string) string {
	return "<div class=\"row\">" + str + "</div>"
}

func title(str string) string {
	return surroundWithColumn(surroundWithH3(str), "ten")
}

func button(siteID string) string {
	var button string = ""
	button += "<form action=\"/editsite/" + siteID + "\" method=\"POST\">"
	button += "<input type=\"submit\" value=\"Edit\">"
	button += "</form>"
	return surroundWithColumn(button, "two")
}

func getDescriptionList(s *site.Entity) string {
	var data string = ""
	data += dataMarkup("Domain", s.Domain)
	data += dataMarkup("Client ID", s.ClientId)
	data += dataMarkup("Private Key", s.PrivateKey)
	data += dataMarkup("Callback URL", s.CallbackURL)
	data += dataMarkup("Cancel URL", s.CancelURL)
	return surroundWithDl(data)
}
