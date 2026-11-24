package elements

import constellation "constellation/src"

func Button(element constellation.WebComponent) string {
	return render("button", element)
}

func Datalist(element constellation.WebComponent) string {
	return render("datalist", element)
}

func Fieldset(element constellation.WebComponent) string {
	return render("fieldset", element)
}

func Form(element constellation.WebComponent) string {
	return render("form", element)
}

func Input(element constellation.WebComponent) string {
	return render("input", element)
}

func Label(element constellation.WebComponent) string {
	return render("label", element)
}

func Legend(element constellation.WebComponent) string {
	return render("legend", element)
}

func Meter(element constellation.WebComponent) string {
	return render("meter", element)
}

func Optgroup(element constellation.WebComponent) string {
	return render("optgroup", element)
}

func Option(element constellation.WebComponent) string {
	return render("option", element)
}

func Output(element constellation.WebComponent) string {
	return render("output", element)
}

func Progress(element constellation.WebComponent) string {
	return render("progress", element)
}

func Select(element constellation.WebComponent) string {
	return render("select", element)
}

func Textarea(element constellation.WebComponent) string {
	return render("textarea", element)
}
