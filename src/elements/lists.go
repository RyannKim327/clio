package elements

import constellation "constellation/src"

func Dd(element constellation.WebComponent) string {
	return render("dd", element)
}

func Dl(element constellation.WebComponent) string {
	return render("dl", element)
}

func Dt(element constellation.WebComponent) string {
	return render("dt", element)
}

func Li(element constellation.WebComponent) string {
	return render("li", element)
}

func Menu(element constellation.WebComponent) string {
	return render("menu", element)
}

func Ol(element constellation.WebComponent) string {
	return render("ol", element)
}

func Ul(element constellation.WebComponent) string {
	return render("ul", element)
}
