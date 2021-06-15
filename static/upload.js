var urlElement = document.querySelector("input");
var scanButton = document.querySelector("button");

function reducer(curr, key) {
	curr[key] = document.querySelector('[name="' + key + '"]')
	return curr
}

var elements = ["url", "title", "author", "tags", "type", "comment"].reduce(reducer, {})

scanButton.addEventListener("click", function(){
	var url = urlElement.value;

	if (url.length < 1 || !url.startsWith("http")) return;

	var req = new XMLHttpRequest();
	req.addEventListener("load", function(){
		if (this.status === 200) {
			var data = JSON.parse(this.responseText)
			for (var key in data) {
				try {
					elements[key].value = data[key]
				} catch(e) {
					console.log(e, key)
				}
			}
		}
	});
	req.open("GET", "/scan");
	req.setRequestHeader("Url", url)
	req.send();
})

