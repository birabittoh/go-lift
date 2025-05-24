// Check for user preference
document.addEventListener("DOMContentLoaded", () => {
  // Check for saved theme preference or use system preference
  const savedTheme = localStorage.getItem("theme");
  if (savedTheme) {
    document.documentElement.setAttribute("data-theme", savedTheme);
    updateThemeToggle(savedTheme);
  } else if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    document.documentElement.setAttribute("data-theme", "dark");
    updateThemeToggle("dark");
  }

  // Listen for system preference changes
  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", (e) => {
      const newTheme = e.matches ? "dark" : "light";
      if (!localStorage.getItem("theme")) {
        document.documentElement.setAttribute("data-theme", newTheme);
        updateThemeToggle(newTheme);
      }
    });

  // Theme toggle button
  const themeToggle = document.getElementById("themeToggle");
  themeToggle.addEventListener("click", () => {
    const currentTheme =
      document.documentElement.getAttribute("data-theme") || "light";
    const newTheme = currentTheme === "light" ? "dark" : "light";

    document.documentElement.setAttribute("data-theme", newTheme);
    localStorage.setItem("theme", newTheme);
    updateThemeToggle(newTheme);
  });

  setTimeout(() => {
    for (item of document.getElementsByClassName("loading")) {
      console.log(item);
      item.classList.remove("loading");
    }
  }, 500);
});

function updateThemeToggle(theme) {
  const themeToggle = document.getElementById("themeToggle");
  themeToggle.textContent = theme === "dark" ? "☀️" : "🌑";
}
