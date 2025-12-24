const API = "http://localhost:8080/api";
let token = localStorage.getItem("token");
let role = localStorage.getItem("role");
let username = localStorage.getItem("user");

// =======================
// INIT
// =======================
window.onload = () => {
  checkAuth();
  loadMangas();
};

// =======================
// AUTH
// =======================
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
    // Kiểm tra sự tồn tại của tất cả các phần tử
    if (token) {
      guest.style.display = "none";
      user.style.display = "inline";
      name.innerText = "Hi, " + username;

      if (role === "admin" && adminBtn) {
        adminBtn.style.display = "inline-block";
      }
    } else {
      guest.style.display = "inline";
      user.style.display = "none";
    }
  } else {
    console.error("❌ Một hoặc nhiều phần tử không tồn tại trong DOM.");
  }
}

// =======================
// MANGA LOGIC
// =======================
let currentMangas = [];

async function loadMangas() {
  try {
    const res = await fetch(API + "/mangas");
    const data = await res.json();

    if (!Array.isArray(data)) {
      console.error("❌ API không trả về mảng", data);
      return;
    }

    currentMangas = data;
    renderMangas(currentMangas);
  } catch (e) {
    console.error("❌ Không load được mangas", e);
  }
}

function renderMangas(list) {
  const grid = document.getElementById("manga-grid");

  if (!grid) {
    console.error("❌ Không tìm thấy #manga-grid trong index.html");
    return;
  }

  if (list.length === 0) {
    grid.innerHTML = "<p style='color:white'>Chưa có manga</p>";
    return;
  }

  grid.innerHTML = list
    .map(
      (m) => `
        <div class="manga-card" onclick="openManga('${m.slug}')">
            <img src="${m.cover}" class="manga-cover"
                 onerror="this.src='https://via.placeholder.com/200x300'">
            <div class="manga-info">
                <div class="manga-title">${m.title}</div>
                <small style="color:#aaa">${m.category}</small>
            </div>
        </div>
    `
    )
    .join("");
}

function searchManga(keyword) {
  const filtered = currentMangas.filter((m) =>
    m.title.toLowerCase().includes(keyword.toLowerCase())
  );
  renderMangas(filtered);
}

// =======================
// DEMO CLICK
// =======================
async function openManga(slug) {
    try {
        const res = await fetch(`${API}/manga/${slug}`);
        const data = await res.json();
        
        const grid = document.getElementById("manga-grid");
        grid.innerHTML = `
            <div style="grid-column: 1/-1; background: #161b40; padding: 20px; border-radius: 10px;">
                <button onclick="loadMangas()" class="btn" style="background: #6c757d; margin-bottom: 20px;">← Quay lại</button>
                <h2 style="color: #ff0055;">${data.title}</h2>
                
                <button onclick="addToWishlist('${slug}')" class="btn" style="background: #f39c12; margin-bottom: 20px;">❤ Thêm vào Wishlist</button>

                <h3 style="border-bottom: 1px solid #333; padding-bottom: 10px;">Danh sách Chapter</h3>
                <div style="display: flex; flex-wrap: wrap; gap: 10px; margin-top: 15px;">
                    ${data.chapters.map(ch => `
                        <button class="btn" style="background: #1f244d; border: 1px solid #ff0055" onclick="readChapter('${slug}', '${ch}')">
                            ${ch.toUpperCase()}
                        </button>
                    `).join('')}
                </div>

                <div style="margin-top: 40px; background: #0b0c2a; padding: 15px; border-radius: 8px;">
                    <h4>Gửi tin nhắn / Comment cho Admin</h4>
                    <textarea id="comment-text" style="width: 100%; background: #1f244d; color: white; border: none; padding: 10px; border-radius: 5px; margin: 10px 0;" rows="3" placeholder="Nhập nội dung..."></textarea>
                    <button onclick="sendMsg()" class="btn" style="background: #ff0055">Gửi tin nhắn</button>
                </div>
            </div>
        `;
    } catch (e) {
        console.error("Lỗi tải chi tiết truyện", e);
    }
}

// Hàm đọc chapter (hiện ảnh)
async function readChapter(slug, chap) {
    const res = await fetch(`${API}/read/${slug}/${chap}`);
    const images = await res.json();
    const grid = document.getElementById("manga-grid");
    
    grid.innerHTML = `
        <div style="grid-column: 1/-1; text-align: center;">
            <button onclick="openManga('${slug}')" class="btn" style="background: #6c757d; margin-bottom: 20px;">← Quay lại danh sách Chap</button>
            <div id="image-container">
                ${images.map(img => `<img src="${img}" style="width: 100%; max-width: 800px; display: block; margin: 0 auto 10px;">`).join('')}
            </div>
        </div>
    `;
}

// Hàm thêm Wishlist
async function addToWishlist(slug) {
    const user = localStorage.getItem("user");
    if (!user) return alert("Vui lòng đăng nhập!");
    const res = await fetch(`${API}/wishlist`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: user, slug: slug })
    });
    const data = await res.json();
    alert(data.message || data.error);
}

// Hàm gửi tin nhắn
async function sendMsg() {
    const user = localStorage.getItem("user");
    const content = document.getElementById("comment-text").value;
    if (!user || !content) return alert("Vui lòng nhập nội dung!");
    
    await fetch(`${API}/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: user, content: content })
    });
    alert("Đã gửi tin nhắn cho Admin!");
    document.getElementById("comment-text").value = "";
}

