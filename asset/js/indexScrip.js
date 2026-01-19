// Toggle class class active untuk menu
const navbarNav = document.querySelector(".navbar-nav");
// ketika menu di click
document.querySelector("#menu").onclick = (e) => {
  navbarNav.classList.toggle("active");
  e.preventDefault();
};

//  dropdown user active
function toggleDropdown() {
  const dropdown = document.getElementById("userDropdown");
  dropdown.style.display =
    dropdown.style.display === "block" ? "none" : "block";
}

// Toggle class active untuk search form
const searchForm = document.querySelector(".search-form");
const searchBox = document.querySelector("#search-box");
const searchButton = document.querySelector("#search-button");

if (searchForm && searchBox && searchButton) {
  searchButton.onclick = (e) => {
    searchForm.classList.toggle("active");
    searchBox.focus();
    e.preventDefault();
  };
}

// // Toggle class active untuk shopping cart
// const shoppingCart = document.querySelector(".shopping-cart");
// document.querySelector("#shopping-cart-button").onclick = (e) => {
//   shoppingCart.classList.toggle("active");
//   e.preventDefault();
// };

// click di liwar elmen untuk menghilangkan nav
const menu = document.querySelector("#menu");
const sb = document.querySelector("#search-button");
// const sc = document.querySelector("#shopping-cart-button");

document.addEventListener("click", function (e) {
  if (!menu.contains(e.target) && !navbarNav.contains(e.target)) {
    navbarNav.classList.remove("active");
  }

  if (sb) {
    if (!sb.contains(e.target) && !searchForm.contains(e.target)) {
      searchForm.classList.remove("active");
    }
  }

  // if (!sc.contains(e.target) && !shoppingCart.contains(e.target)) {
  //   shoppingCart.classList.remove("active");
  // }
});

// Opsional: tutup dropdown jika klik di luar
document.addEventListener("click", function (e) {
  const userLink = document.querySelector(".navbar-user");
  const dropdown = document.getElementById("userDropdown");
  if (!userLink.contains(e.target) && !dropdown.contains(e.target)) {
    dropdown.style.display = "none";
  }
});

//tampilkan priviw gambar yang ada di add produk
const profileImageInput = document.getElementById("profileImageInput");

if (profileImageInput) {
  profileImageInput.addEventListener("change", function (event) {
    const file = event.target.files[0];
    if (file) {
      const preview = document.getElementById("profilePicture");
      if (preview) {
        preview.src = URL.createObjectURL(file);
      }
    }
  });
}

window.addEventListener("DOMContentLoaded", async () => {
  const res = await fetch("/produk/cartCount");
  const data = await res.json();
  document.getElementById("cart-count").textContent = data.cart_count;
});

async function increaseItem(id) {
  const res = await fetch(`/cart/increase/${id}`, { method: "POST" });
  const result = await res.json();

  if (res.ok) {
    updateUI(id, result.itemQuantity, result.itemSubtotal, result.totalCart);
    const cartCountEl = document.getElementById("cart-count");
    cartCountEl.textContent = result.cart_count;
  }
}

async function decreaseItem(id) {
  const res = await fetch(`/cart/decrease/${id}`, { method: "POST" });
  const result = await res.json();

  if (res.ok) {
    if (result.removed) {
      document.getElementById(`row-${id}`).remove();
    } else {
      updateUI(id, result.itemQuantity, result.itemSubtotal, result.totalCart);
    }
    const cartCountEl = document.getElementById("cart-count");
    cartCountEl.textContent = result.cart_count;
  }
}

async function deleteItem(id) {
  const res = await fetch(`/cart/decrease/${id}`, { method: "POST" });
  const result = await res.json();

  if (res.ok) {
    document.getElementById(`row-${id}`).remove();
    document.getElementById("total").textContent = `Rp ${result.totalCart}`;
    const cartCountEl = document.getElementById("cart-count");
    cartCountEl.textContent = result.cart_count;
  }
}

function updateUI(id, qty, subtotal, total) {
  document.getElementById(`qty-${id}`).textContent = qty;
  document.getElementById(`subtotal-${id}`).textContent = `Rp ${subtotal}`;
  document.getElementById("total").textContent = `Rp ${total}`;
}
