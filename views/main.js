const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

signUpButton.addEventListener('click', () => {
	container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
	container.classList.remove("right-panel-active");
});

// strong password filtering with regex 

    // Function to check if password meets the requirements
    function isPasswordValid(password) {
        var regex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()\-_=+\\\|\[\]{};:'",.<>\/?]).{8,}$/;
        return regex.test(password);
    }

    // Function to enable/disable signup button based on password validity
    function updateSignupButton() {
        var password = document.getElementById("password").value;
        var signupBtn = document.getElementById("signup-btn");

        if (isPasswordValid(password)) {
            signupBtn.disabled = false; // Enable button if password is valid
        } else {
            signupBtn.disabled = true; // Disable button if password is invalid
        }
    }

    // Add event listeners to password input for live validation
    document.getElementById("password").addEventListener("input", updateSignupButton);

    // Add event listener to form submission for final validation
    document.getElementById("signup-form").addEventListener("submit", function(event) {
        var password = document.getElementById("password").value;

        if (!isPasswordValid(password)) {
            event.preventDefault(); // Prevent form submission if password is invalid
			document.getElementById("error-message").innerText = "Password must be at least 8 characters long, contain at least one uppercase letter, one digit, and one special character.";
        }
    });

	document.getElementById("password").addEventListener("input", function(event) {
		var password = event.target.value;
		var errorMessage = document.getElementById("password-error");
	
		if (!isPasswordValid(password)) {
			errorMessage.innerText = "Password must be at least 8 characters long, contain at least one uppercase letter, number and one special character.";
		} else {
			errorMessage.innerText = "";
		}
	});


// Function to get query string parameter by name
function getQueryStringParameter(name) {
	const urlParams = new URLSearchParams(window.location.search);
	return urlParams.get(name);
}

// Function to display message if exists in query string
function displayMessage() {
    const message = getQueryStringParameter('message');
    const type = getQueryStringParameter('type');
    let icon = 'success'; // Default icon
    if (type === 'error') {
        icon = 'error'; // Change icon to error if type is 'error'
    }

    if (message) {
        // Display message using Sweet Alert
        Swal.fire({
            icon: icon,
            title: type === 'error' ? 'Error!' : 'Success!',
            text: message,
            showConfirmButton: false,
            timer: 2000
        });
    }
}

window.onload = displayMessage;



