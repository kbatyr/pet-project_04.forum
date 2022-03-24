var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

var expanded = false;

function showCheckboxes() {
	var checkboxes = document.getElementById("checkboxes");
	if (!expanded) {
		checkboxes.style.display = "block";
		expanded = true;
	} else {
		checkboxes.style.display = "none";
		expanded = false;
	}
}
// like & dislike btn
var btn1 = document.querySelector('#green');
var btn2 = document.querySelector('#red');

// console.log(btn1, btn2)

btn1.addEventListener('click', function () {
	console.log('here')

	if (btn2.classList.contains('red')) {
		btn2.classList.remove('red');
	}
	this.classList.toggle('green');
});

btn2.addEventListener('click', function () {

	if (btn1.classList.contains('green')) {
		btn1.classList.remove('green');
	}
	this.classList.toggle('red');
});