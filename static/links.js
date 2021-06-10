var links = document.querySelectorAll("a")
var details = document.getElementById("details")

function getData(element, attribute) {
	return element.getAttribute("data-" + attribute);
}

links.forEach(function(link){
	link.addEventListener("click", function(e){
		e.preventDefault();
		var url = link.href,
			title = link.innerText,
			date = getData(link, "date"),
			comment = getData(link, "comment"),
			tags = getData(link, "tags"),
			author = getData(link, "author");

		var linkElement = document.createElement("a"),
			tagsElement = document.createElement("tags"),
			commentElement = document.createElement("p")

		linkElement.className = "link"
		linkElement.href = url
		linkElement.innerText = title + " ↗"
		linkElement.target = "_blank"
		linkElement.rel = "noopener noreferrer"

		tagsElement.className = "tags"
		tagsElement.innerText = [tags, date, author].join(" — ")

		commentElement.className = "comment"
		commentElement.innerText = comment

		details.innerText = ""
		details.appendChild(linkElement)
		details.appendChild(tagsElement)
		details.appendChild(commentElement)

		details.style.display = "block"
	})
})

