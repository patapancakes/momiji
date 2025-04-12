var allowedVars = ["accent", "border", "bg-odd", "bg-even", "text", "msg-text", "link"];
window.addEventListener("message", function(e) {
	if (e.data["var"] == "mode") {
		if (e.data["val"] == "dark") {
			document.documentElement.classList.add("dark");
		} else {
			document.documentElement.classList.remove("dark");
		}
	} else if (allowedVars.includes(e.data["var"])) {
		if (e.data["val"].match(/^#([0-9a-f]{3}){1,2}$/i)) {
			document.documentElement.style.setProperty("--"+e.data["var"], e.data["val"]);
		} else {
			document.documentElement.style.removeProperty("--"+e.data["var"]);
		}
	}
});
