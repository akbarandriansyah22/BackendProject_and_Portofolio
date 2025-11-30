// ============================================
// SMOOTH SCROLLING FOR NAVIGATION LINKS
// ============================================
document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
  anchor.addEventListener("click", function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute("href"));

    if (target) {
      const offsetTop = target.offsetTop - 70; // Offset untuk navbar
      window.scrollTo({
        top: offsetTop,
        behavior: "smooth",
      });
    }
  });
});

// ============================================
// NAVBAR SCROLL EFFECT & ACTIVE LINK
// ============================================
const navbar = document.querySelector(".navbar");
const sections = document.querySelectorAll("section[id]");
const navLinks = document.querySelectorAll(".navbar nav a");

window.addEventListener("scroll", function () {
  // Navbar scroll effect
  if (window.scrollY > 100) {
    navbar.classList.add("scrolled");
  } else {
    navbar.classList.remove("scrolled");
  }

  // Active navigation highlight
  let scrollPosition = window.scrollY + 150;

  sections.forEach((section) => {
    const sectionTop = section.offsetTop;
    const sectionHeight = section.offsetHeight;
    const sectionId = section.getAttribute("id");

    if (
      scrollPosition >= sectionTop &&
      scrollPosition < sectionTop + sectionHeight
    ) {
      navLinks.forEach((link) => {
        link.classList.remove("active");
        if (link.getAttribute("href") === `#${sectionId}`) {
          link.classList.add("active");
        }
      });
    }
  });
});

// ============================================
// TYPING EFFECT FOR HERO SUBTITLE
// ============================================
const typedTextSpan = document.querySelector(".typed-text");
const cursor = document.querySelector(".cursor");

if (typedTextSpan) {
  const textArray = [
    "DevOps Enthusiast",
    "Git & GitHub Explorer",
    "Version Control Learner",
    "Future DevOps Engineer",
  ];

  const typingDelay = 100;
  const erasingDelay = 50;
  const newTextDelay = 2000;
  let textArrayIndex = 0;
  let charIndex = 0;

  function type() {
    if (charIndex < textArray[textArrayIndex].length) {
      typedTextSpan.textContent += textArray[textArrayIndex].charAt(charIndex);
      charIndex++;
      setTimeout(type, typingDelay);
    } else {
      setTimeout(erase, newTextDelay);
    }
  }

  function erase() {
    if (charIndex > 0) {
      typedTextSpan.textContent = textArray[textArrayIndex].substring(
        0,
        charIndex - 1
      );
      charIndex--;
      setTimeout(erase, erasingDelay);
    } else {
      textArrayIndex++;
      if (textArrayIndex >= textArray.length) textArrayIndex = 0;
      setTimeout(type, typingDelay + 500);
    }
  }

  // Start typing effect when page loads
  window.addEventListener("load", function () {
    setTimeout(type, newTextDelay);
  });
}

// ============================================
// SCROLL ANIMATIONS (FADE IN EFFECT)
// ============================================
const observerOptions = {
  threshold: 0.15,
  rootMargin: "0px 0px -50px 0px",
};

const observer = new IntersectionObserver(function (entries) {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      entry.target.classList.add("appear");
    }
  });
}, observerOptions);

// Apply fade-in animation to sections and cards
document.addEventListener("DOMContentLoaded", function () {
  // Animate sections
  const animatedSections = document.querySelectorAll("section:not(.hero)");
  animatedSections.forEach((section) => {
    section.classList.add("fade-in");
    observer.observe(section);
  });

  // Animate cards with stagger effect
  const cards = document.querySelectorAll(
    ".project-card, .skill-category, .learning-card"
  );
  cards.forEach((card, index) => {
    card.classList.add("fade-in");
    card.style.transitionDelay = `${index * 0.1}s`;
    observer.observe(card);
  });
});

// ============================================
// BACK TO TOP BUTTON
// ============================================
const backToTopButton = document.getElementById("backToTop");

if (backToTopButton) {
  window.addEventListener("scroll", function () {
    if (window.scrollY > 300) {
      backToTopButton.classList.add("visible");
    } else {
      backToTopButton.classList.remove("visible");
    }
  });

  backToTopButton.addEventListener("click", function () {
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  });
}

// ============================================
// PROJECT CARD INTERACTIONS
// ============================================
const projectCards = document.querySelectorAll(".project-card");

projectCards.forEach((card) => {
  // Add click effect
  card.addEventListener("click", function (e) {
    // Don't trigger if clicking on a link
    if (e.target.tagName !== "A") {
      const projectTitle = this.querySelector("h4").textContent;
      console.log(`Project clicked: ${projectTitle}`);
    }
  });

  // Add keyboard navigation
  card.setAttribute("tabindex", "0");
  card.addEventListener("keypress", function (e) {
    if (e.key === "Enter") {
      const link = this.querySelector(".project-link");
      if (link) link.click();
    }
  });
});

// ============================================
// SKILLS CATEGORY EXPAND EFFECT
// ============================================
const skillCategories = document.querySelectorAll(".skill-category");

skillCategories.forEach((category) => {
  category.addEventListener("mouseenter", function () {
    const items = this.querySelectorAll(".skill-item");
    items.forEach((item, index) => {
      setTimeout(() => {
        item.style.transform = "scale(1.05)";
      }, index * 50);
    });
  });

  category.addEventListener("mouseleave", function () {
    const items = this.querySelectorAll(".skill-item");
    items.forEach((item) => {
      item.style.transform = "scale(1)";
    });
  });
});

// ============================================
// CONSOLE WELCOME MESSAGE
// ============================================
console.log(
  "%c🚀 Welcome to Akbar's Portfolio!",
  "color: #2563eb; font-size: 20px; font-weight: bold;"
);
console.log(
  "%c👨‍💻 Currently learning: Git, GitHub, GitLab",
  "color: #10b981; font-size: 14px;"
);
console.log(
  "%c📧 Contact: akbarandriansyah73@gmail.com",
  "color: #64748b; font-size: 12px;"
);
console.log(
  "%c💡 Tip: Check out my GitHub for more projects!",
  "color: #f59e0b; font-size: 12px;"
);

// ============================================
// PERFORMANCE OPTIMIZATION
// ============================================
// Lazy loading for images (if you add images later)
if ("IntersectionObserver" in window) {
  const imageObserver = new IntersectionObserver((entries, observer) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        const img = entry.target;
        if (img.dataset.src) {
          img.src = img.dataset.src;
          img.classList.remove("lazy");
          observer.unobserve(img);
        }
      }
    });
  });

  const lazyImages = document.querySelectorAll("img.lazy");
  lazyImages.forEach((img) => imageObserver.observe(img));
}

// ============================================
// PREVENT ANIMATION ON PAGE LOAD
// ============================================
window.addEventListener("load", function () {
  document.body.classList.add("loaded");
});

// ============================================
// KEYBOARD NAVIGATION IMPROVEMENTS
// ============================================
document.addEventListener("keydown", function (e) {
  // Press 'H' to go to home
  if (e.key === "h" || e.key === "H") {
    if (!e.target.matches("input, textarea")) {
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
  }
});

// ============================================
// EMAIL COPY TO CLIPBOARD (BONUS FEATURE)
// ============================================
const emailElements = document.querySelectorAll(".contact-item p");

emailElements.forEach((element) => {
  const text = element.textContent.trim();
  if (text.includes("@")) {
    element.style.cursor = "pointer";
    element.title = "Click to copy email";

    element.addEventListener("click", function (e) {
      e.preventDefault();
      const email = this.textContent.trim();

      // Copy to clipboard
      if (navigator.clipboard) {
        navigator.clipboard
          .writeText(email)
          .then(() => {
            // Show copied notification
            const originalText = this.textContent;
            this.textContent = "✓ Email copied!";
            this.style.color = "#10b981";

            setTimeout(() => {
              this.textContent = originalText;
              this.style.color = "";
            }, 2000);
          })
          .catch((err) => {
            console.error("Failed to copy:", err);
          });
      }
    });
  }
});

// ============================================
// ROADMAP ITEM ANIMATION ON CLICK
// ============================================
const roadmapItems = document.querySelectorAll(".roadmap li");

roadmapItems.forEach((item) => {
  item.addEventListener("click", function () {
    this.style.transform = "translateX(15px) scale(1.02)";
    setTimeout(() => {
      this.style.transform = "";
    }, 300);
  });
});

// ============================================
// PREVENT EMPTY LINKS
// ============================================
document.querySelectorAll('a[href="#"]').forEach((link) => {
  link.addEventListener("click", function (e) {
    e.preventDefault();
  });
});
