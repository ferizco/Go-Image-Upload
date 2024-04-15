
    // Dapatkan elemen tombol dan modal
    var openModalBtn = document.getElementById("open-modal-btn");
    var modal = document.getElementById("myModal");

    // Ketika tombol ditekan, tampilkan modal
    openModalBtn.addEventListener("click", function() {
        modal.classList.add("show");
    });

    // Dapatkan elemen tombol untuk menutup modal
    var closeModalBtn = document.querySelector(".close");

    // Ketika tombol untuk menutup modal ditekan, sembunyikan modal
    closeModalBtn.addEventListener("click", function() {
        modal.classList.remove("show");
    });

    // Dapatkan form
    var submitForm = document.getElementById("submit-form");

 
    var slideIndex = 1;
    showSlides(slideIndex);
    
    function plusSlides(n) {
        showSlides(slideIndex += n);
    }
    
    function showSlides(n) {
        var i;
        var slides = document.getElementsByClassName("mySlides");
        if (slides.length === 0) return; // Tidak ada slide, keluar dari fungsi
        if (n > slides.length) {slideIndex = 1}    
        if (n < 1) {slideIndex = slides.length}
        for (i = 0; i < slides.length; i++) {
            slides[i].style.display = "none";  
        }
        // Menampilkan 2 gambar
        if (slides[slideIndex - 1]) {
            slides[slideIndex - 1].style.display = "block";
        }
    
        // Menampilkan gambar berikutnya jika ada
        if (slides[slideIndex]) {
            slides[slideIndex].style.display = "block";
        } else {
            // Jika slides[slideIndex] adalah undefined, kembalikan ke slides pertama
            slides[0].style.display = "block";
        }
    }


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
    } console.log(message,type)
}

window.onload = displayMessage;
    
