/**
 * @function
 * An onclick event listener to toggle the comment input form.
 * By default the input box is hidden and this will how it
 */
function showForm() {
	const x = document.getElementById("commentForm");
	if (x.style.display === "flex") {
		x.style.display = "none";
	} else {
		x.style.display = "flex";
	}
}

const commentListener = () => {
	const a = document.getElementById("addComment");
	a.onclick = showForm;
};

commentListener();
