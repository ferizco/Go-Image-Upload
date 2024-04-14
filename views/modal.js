
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
        if (n > slides.length) {slideIndex = 1}    
        if (n < 1) {slideIndex = slides.length}
        for (i = 0; i < slides.length; i++) {
            slides[i].style.display = "none";  
        }
        // Menampilkan 2 gambar
        slides[slideIndex - 1].style.display = "block";
        if (slides[slideIndex]) {
            slides[slideIndex].style.display = "block";
        }
    }
    
