/**
 * @function
 * used as an onclick event listerner to show the image within a
 * modal as a bigger full screen version
 */
function enlargeModal() {
	const m = document.getElementById("bigImg");
	m.style.backgroundImage = `url('${this.src}')`;
	m.showModal();
}

const addImageListeners = () => {
	const imgs = document.querySelectorAll("img");
	if (imgs.length === 0) {
		return;
	}

	for (let k = 0; k < imgs.length; k++) {
		imgs[k].onclick = enlargeModal;
	}

	const b = document.querySelector("#bigImg button");
	b.onclick = function () {
		const d = document.getElementById("bigImg");
		d.close();
	};
};

addImageListeners();
