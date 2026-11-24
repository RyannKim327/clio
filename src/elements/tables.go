package elements

import constellation "constellation/src"

func Caption(element constellation.WebComponent) string {
	return render("caption", element)
}

func Col(element constellation.WebComponent) string {
	return render("col", element)
}

func Colgroup(element constellation.WebComponent) string {
	return render("colgroup", element)
}

func Table(element constellation.WebComponent) string {
	return render("table", element)
}

func Tbody(element constellation.WebComponent) string {
	return render("tbody", element)
}

func Td(element constellation.WebComponent) string {
	return render("td", element)
}

func Tfoot(element constellation.WebComponent) string {
	return render("tfoot", element)
}

func Th(element constellation.WebComponent) string {
	return render("th", element)
}

func Thead(element constellation.WebComponent) string {
	return render("thead", element)
}

func Tr(element constellation.WebComponent) string {
	return render("tr", element)
}
