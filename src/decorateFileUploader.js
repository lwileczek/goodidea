/**
 * @function
 * add little dots next to the upload button when users select images
 * to let the user know we have their files prepared to go to the server
 */
const setUploaderDecoration = () => {
    document.getElementById("newImgUpload").onchange = function() {
        let content = "";
        if (this.files.length > 0) {
            const m = this.files.length > 3 ? 3 : this.files.length;
            for (let d = 0; d < m; d++) {
                content += " &#9679;"
            }
        }
        document.getElementById("uploadFile").innerHTML = content;
    };
};

setUploaderDecoration();
