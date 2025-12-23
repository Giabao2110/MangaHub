const API = "http://localhost:8080/api";
let token = localStorage.getItem("token");
let role = localStorage.getItem("role");
let username = localStorage.getItem("user");

// --- KHỞI TẠO ---
window.onload = () => {
    loadMangas();
    checkAuth();
};

// --- LOGIC AUTHENTICATION (QUAN TRỌNG) ---

// 1. Hàm Đăng Ký
async function register() {
    const u = prompt("Tạo tên đăng nhập mới:");
    if (!u) return;
    const p = prompt("Tạo mật khẩu:");
    if (!p) return;

    try {
        const res = await fetch(API + "/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username: u, password: p })
        });
        const data = await res.json();
        
        if (res.ok) {
            alert("✅ Đăng ký thành công! Bạn có thể đăng nhập ngay.");
        } else {
            alert("❌ Lỗi: " + data.error);
        }
    } catch (e) {
        alert("Không thể kết nối Server!");
    }
}

// 2. Hàm Đăng Nhập
async function login() {
    const u = prompt("Nhập Username:");
    if (!u) return;
    const p = prompt("Nhập Password:");
    if (!p) return;

    try {
        const res = await fetch(API + "/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username: u, password: p })
        });
        const data = await res.json();

        if (res.ok) {
            localStorage.setItem("token", data.token);
            localStorage.setItem("role", data.role);
            localStorage.setItem("user", data.username);
            location.reload(); // Tải lại trang để cập nhật giao diện
        } else {
            alert("❌ " + data.error);
        }
    } catch (e) {
        alert("Lỗi kết nối Server!");
    }
}

// 3. Hàm Đăng Xuất
function logout() {
    localStorage.clear();
    location.reload();
}

// 4. Kiểm tra trạng thái đăng nhập để đổi nút
function checkAuth() {
    if (token) {
        // Nếu đã đăng nhập: Ẩn nút Đăng ký/Đăng nhập, hiện nút Thoát
        document.getElementById("guest-area").style.display = "none";
        document.getElementById("user-area").style.display = "inline";
        document.getElementById("username-display").innerText = "Hi, " + username;
        
        if (role === 'admin') {
            document.getElementById("nav-admin").style.display = "inline-block";
        }
    } else {
        // Nếu chưa đăng nhập: Hiện nút Đăng ký/Đăng nhập
        document.getElementById("guest-area").style.display = "inline";
        document.getElementById("user-area").style.display = "none";
    }
}

// --- LOGIC TRUYỆN (Giữ nguyên như cũ) ---
let currentMangas = [];
async function loadMangas() {
    const res = await fetch(API + "/mangas");
    currentMangas = await res.json();
    renderMangas(currentMangas);
}

function renderMangas(list) {
    const grid = document.getElementById("manga-grid");
    grid.innerHTML = list.map(m => `
        <div class="manga-card" onclick="openManga('${m.slug}')">
            <img src="${m.cover}" class="manga-cover" onerror="this.src='https://via.placeholder.com/200x300'">
            <div class="manga-info">
                <div class="manga-title">${m.title}</div>
                <small style="color: #aaa">${m.category}</small>
            </div>
        </div>
    `).join('');
}

function searchManga(keyword) {
    const filtered = currentMangas.filter(m => m.title.toLowerCase().includes(keyword.toLowerCase()));
    renderMangas(filtered);
}