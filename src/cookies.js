/**
 * @Function{setVoteListeners}
 * List over the voting buttons and set listeners to set cookies if a user votes
 * Then the next time they visit the page they cannot recast their vote in the same direction
 */
const setVoteListeners = () => {
	const radios = document.querySelectorAll("#task-list input[type=radio]");
	for (let j = 0; j < radios.length; j++) {
		const id = radios[j].id;
		const voteInfo = id.split("-");
		radios[j].onclick = function () {
			setCookie(voteInfo[0], voteInfo[1], 1);
		};
	}
};

/**
 * @Function{setCookie}
 * record that a user voted on this task by setting a cookie
 * @param{string} taskId - the id of the taskId which is an integer represented as a string
 * @param{string} vote - the direction the user voted, up or down
 * @param{number} exdays - how long the cookie should live in days
 */
const setCookie = (taskId, vote, exdays) => {
	//document.cookie = "username=John Doe; expires=Thu, 18 Dec 2013 12:00:00 UTC";
	const d = new Date();
	d.setTime(d.getTime() + exdays * 24 * 60 * 60 * 1000);
	const expires = `expires=${d.toUTCString()}`;
    document.cookie = `${taskId}=${vote};${expires};path=/`;
};

/**
 * @Function{getCookie}
 * Check if a cookie exists and if so, return the value of that cookie
 * @param{string} cname - the name of the cookie to look for, likely the ID of the task
 * @returns{string} and empty string if not found or the cookie value if it does
 */
const getCookie = (cname) => {
    const name = `${cname}=`;
	const decodedCookie = decodeURIComponent(document.cookie);
	const ca = decodedCookie.split(";");
	for (let i = 0; i < ca.length; i++) {
		let c = ca[i];
		while (c.charAt(0) === " ") {
			c = c.substring(1);
		}
		if (c.indexOf(name) === 0) {
			return c.substring(name.length, c.length);
		}
	}
	return "";
};

function listAllCookieNames() {
	const allCookies = document.cookie;
	const cookieArray = allCookies.split(";");
	const cookieNames = [];
	for (const cookie of cookieArray) {
		const [name, value] = cookie.split("=");
        cookieNames.push(`${name}-${value}`);
	}

	return cookieNames;
}

setVoteListeners();
