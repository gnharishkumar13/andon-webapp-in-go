(function() {
	const elem = document.querySelector(".clock");
	if (!elem) {
		console.error("Cannot find element to attach clock to");
		return;
	}
	updateClock();
	setInterval(updateClock, 1000);

	function updateClock() {
		const d = new Date();
		const hms = [d.getHours(), d.getMinutes(), d.getSeconds()];
		hms.forEach((v, i) => {
			hms[i] = v < 10 ? `0${v}` : `${v}`;
		});
		elem.textContent = hms.join(":");
	}
})();
