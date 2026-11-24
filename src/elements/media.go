package elements

import constellation "constellation/src"

func Area(element constellation.WebComponent) string {
	return render("area", element)
}

func Audio(element constellation.WebComponent) string {
	return render("audio", element)
}

func Canvas(element constellation.WebComponent) string {
	return render("canvas", element)
}

func Embed(element constellation.WebComponent) string {
	return render("embed", element)
}

func Hr(element constellation.WebComponent) string {
	return render("hr", element)
}

func Iframe(element constellation.WebComponent) string {
	return render("iframe", element)
}

func Img(element constellation.WebComponent) string {
	return render("img", element)
}

func Map(element constellation.WebComponent) string {
	return render("map", element)
}

func Object(element constellation.WebComponent) string {
	return render("object", element)
}

func Param(element constellation.WebComponent) string {
	return render("param", element)
}

func Picture(element constellation.WebComponent) string {
	return render("picture", element)
}

func Source(element constellation.WebComponent) string {
	return render("source", element)
}

func Track(element constellation.WebComponent) string {
	return render("track", element)
}

func Video(element constellation.WebComponent) string {
	return render("video", element)
}
