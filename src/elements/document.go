package elements

import constellation "constellation/src"

func Base(element constellation.WebComponent) string {
	return render("base", element)
}

func Link(element constellation.WebComponent) string {
	return render("link", element)
}

func Meta(element constellation.WebComponent) string {
	return render("meta", element)
}

func Noscript(element constellation.WebComponent) string {
	return render("noscript", element)
}

func Script(element constellation.WebComponent) string {
	return render("script", element)
}

func Slot(element constellation.WebComponent) string {
	return render("slot", element)
}

func Style(element constellation.WebComponent) string {
	return render("style", element)
}
