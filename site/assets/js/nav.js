document.addEventListener("DOMContentLoaded", () => {
  const btn = document.querySelector(".nav-toggle");
  const nav = document.querySelector(".site-nav");

  btn.addEventListener("click", () => {
    nav.classList.toggle("open");

    const isOpen = nav.classList.contains("open");
    btn.setAttribute("aria-expanded", isOpen);
  });
});