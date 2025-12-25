const API = "http://localhost:8080/api";
let token = localStorage.getItem("token");
let role = localStorage.getItem("role");
let username = localStorage.getItem("user");

// Pagination state
let currentMangas = [];
let displayedCount = 8;
const increment = 8;

window.onload = () => {
  checkAuth();
  loadMangas();
};

function logout() {
  localStorage.clear();
  location.reload();
}

function checkAuth() {
  const guest = document.getElementById("guest-area");
  const user = document.getElementById("user-area");
  const name = document.getElementById("username-display");
  const adminBtn = document.getElementById("nav-admin");

  if (guest && user && name) {
    if (token) {
      guest.style.display = "none";
      user.style.display = "inline";
      name.innerText = "Hi, " + username;

      if (role === "admin" && adminBtn) adminBtn.style.display = "inline-block";
    } else {
      guest.style.display = "inline";
      user.style.display = "none";
    }
  }
}

// UC-002:
async function loadMangas() {
  try {
    const res = await fetch(API + "/mangas");
    const data = await res.json();
    if (!Array.isArray(data)) return;

    currentMangas = data;
    renderMangas(currentMangas.slice(0, displayedCount));
    updateLoadMoreButton();
  } catch (e) {
    console.error("❌ Failed to load mangas", e);
  }
}

function renderMangas(list, append = false) {
  const grid = document.getElementById("manga-grid");
  if (!grid) return;

  const html = list
    .map(
      (m) => `
        <div class="manga-card" onclick="openManga('${m.slug}')">
            <div class="manga-cover-wrapper">
                <img src="${m.cover}" class="manga-cover" loading="lazy"
                     onerror="this.src='https://via.placeholder.com/200x300?text=No+Cover'">
            </div>
            <div class="manga-info">
                <div class="manga-title">${m.title}</div>
                <small style="color:#aaa; font-weight: 500;">${
                  m.category || "Manga"
                }</small>
            </div>
        </div>
    `
    )
    .join("");

  if (append) {
    grid.insertAdjacentHTML("beforeend", html);
  } else {
    grid.innerHTML = html || "<p>No manga found.</p>";
  }
}

// Tính năng "Xem thêm"
function loadMore() {
  const nextBatch = currentMangas.slice(
    displayedCount,
    displayedCount + increment
  );
  renderMangas(nextBatch, true);
  displayedCount += increment;
  updateLoadMoreButton();
}

function updateLoadMoreButton() {
  let btnContainer = document.getElementById("load-more-container");

  // Nếu chưa có container thì tạo mới
  if (!btnContainer) {
    btnContainer = document.createElement("div");
    btnContainer.id = "load-more-container";
    btnContainer.className = "load-more-container";
    document.getElementById("manga-grid").after(btnContainer);
  }

  if (displayedCount < currentMangas.length) {
    btnContainer.innerHTML = `<button class="btn btn-load-more" onclick="loadMore()">View More</button>`;
  } else {
    btnContainer.innerHTML = ""; // Ẩn nút nếu đã hết truyện
  }
}

function searchManga(keyword) {
  const filtered = currentMangas.filter((m) =>
    m.title.toLowerCase().includes(keyword.toLowerCase())
  );
  renderMangas(filtered);
  // Ẩn nút load more khi đang tìm kiếm để tránh xung đột logic
  const btnContainer = document.getElementById("load-more-container");
  if (btnContainer) btnContainer.innerHTML = "";
}

// UC-004: Xem chi tiết truyện
async function openManga(slug) {
  console.log("Opening manga:", slug);
  try {
    const res = await fetch(`${API}/manga/${slug}`);
    const data = await res.json();

    const grid = document.getElementById("manga-grid");
    // Ẩn nút "View More" khi vào xem chi tiết
    const btnContainer = document.getElementById("load-more-container");
    if (btnContainer) btnContainer.innerHTML = "";

    grid.innerHTML = `
            <div class="manga-detail-view" style="grid-column: 1/-1; background: #161b40; padding: 20px; border-radius: 10px;">
                <button onclick="location.reload()" class="btn" style="background: #6c757d; margin-bottom: 20px;">← Back to List</button>
                <h2 style="color: #ff0055; margin-bottom: 15px;">${
                  data.title
                }</h2>
                
                <button onclick="addToWishlist('${slug}')" class="btn" style="background: #f39c12; margin-bottom: 25px;">❤ Add to Wishlist</button>

                <h3 style="border-bottom: 1px solid #333; padding-bottom: 10px; margin-bottom: 15px;">Chapter List</h3>
                <div style="display: flex; flex-wrap: wrap; gap: 10px;">
                    ${data.chapters
                      .map(
                        (ch) => `
                        <button class="btn" style="background: #1f244d; border: 1px solid #ff0055; min-width: 100px;" 
                                onclick="readChapter('${slug}', '${ch}')">
                            ${ch.toUpperCase()}
                        </button>
                    `
                      )
                      .join("")}
                </div>
                <div id="reader-area" style="margin-top: 30px;"></div>
            </div>
        `;
  } catch (e) {
    console.error("❌ Error loading manga details:", e);
    alert("Could not load manga details!");
  }
}

// UC-005: Đọc Chapter
// Hàm đọc chapter (hiện ảnh)
async function readChapter(slug, chap) {
  const res = await fetch(`${API}/read/${slug}/${chap}`);
  const images = await res.json();
  const grid = document.getElementById("manga-grid");

  grid.innerHTML = `
        <div style="grid-column: 1/-1; text-align: center;">
            <button onclick="openManga('${slug}')" class="btn" style="background: #6c757d; margin-bottom: 20px;">← Quay lại danh sách Chap</button>
            <div id="image-container">
                ${images
                  .map(
                    (img) =>
                      `<img src="${img}" style="width: 100%; max-width: 800px; display: block; margin: 0 auto 10px;">`
                  )
                  .join("")}
            </div>
        </div>
    `;
}

// UC-006: Thêm vào danh sách yêu thích
// Hàm thêm Wishlist
async function addToWishlist(slug) {
  const user = localStorage.getItem("user");
  if (!user) return alert("Vui lòng đăng nhập!");
  const res = await fetch(`${API}/wishlist`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: user, slug: slug }),
  });
  const data = await res.json();
  alert(data.message || data.error);
}

// Hàm gửi tin nhắn
async function sendMsg() {
  const user = localStorage.getItem("user");
  const content = document.getElementById("comment-text").value;
  if (!user || !content) return alert("Vui lòng nhập nội dung!");

  const response = await fetch(`${API}/messages`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: user, content: content }),
  });

  if (response.ok) {
    alert("Đã gửi tin nhắn cho Admin!");

    // THÊM DÒNG NÀY: Tự hiển thị tin nhắn vừa gửi vào vùng hiển thị trên web
    const msgList = document.getElementById("sent-messages-list");
    if (msgList) {
      const time = new Date().toLocaleTimeString();
      msgList.insertAdjacentHTML(
        "afterbegin",
        `
            <div class="msg-item">
                <strong>Bạn:</strong> ${content} <small>(${time})</small>
            </div>
        `
      );
    }

    document.getElementById("comment-text").value = "";
  }
}
