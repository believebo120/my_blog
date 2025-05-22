document.addEventListener('DOMContentLoaded', function() {
  // 登录表单验证
  const loginForm = document.getElementById('loginForm');
  if (loginForm) {
      loginForm.addEventListener('submit', function(e) {
          e.preventDefault();
          
          const username = document.getElementById('username').value.trim();
          const password = document.getElementById('password').value.trim();
          let isValid = true;
          
          // 验证用户名
          if (username === '') {
              document.getElementById('usernameError').textContent = '请输入用户名';
              isValid = false;
          } else {
              document.getElementById('usernameError').textContent = '';
          }
          
          // 验证密码
          if (password === '') {
              document.getElementById('passwordError').textContent = '请输入密码';
              isValid = false;
          } else if (password.length < 6) {
              document.getElementById('passwordError').textContent = '密码长度不能少于6位';
              isValid = false;
          } else {
              document.getElementById('passwordError').textContent = '';
          }
          
          if (isValid) {
              alert('登录成功！');
              window.location.href = 'blog.html';
          }
      });
  }
  
  // 注册表单验证
  const registerForm = document.getElementById('registerForm');
  if (registerForm) {
      registerForm.addEventListener('submit', function(e) {
          e.preventDefault();
          
          const username = document.getElementById('regUsername').value.trim();
          const email = document.getElementById('regEmail').value.trim();
          const password = document.getElementById('regPassword').value.trim();
          const confirmPassword = document.getElementById('regConfirmPassword').value.trim();
          let isValid = true;
          
          // 验证用户名
          if (username === '') {
              document.getElementById('regUsernameError').textContent = '请输入用户名';
              isValid = false;
          } else if (username.length < 4) {
              document.getElementById('regUsernameError').textContent = '用户名长度不能少于4位';
              isValid = false;
          } else {
              document.getElementById('regUsernameError').textContent = '';
          }
          
          // 验证邮箱
          if (email === '') {
              document.getElementById('regEmailError').textContent = '请输入邮箱';
              isValid = false;
          } else if (!validateEmail(email)) {
              document.getElementById('regEmailError').textContent = '请输入有效的邮箱地址';
              isValid = false;
          } else {
              document.getElementById('regEmailError').textContent = '';
          }
          
          // 验证密码
          if (password === '') {
              document.getElementById('regPasswordError').textContent = '请输入密码';
              isValid = false;
          } else if (password.length < 6) {
              document.getElementById('regPasswordError').textContent = '密码长度不能少于6位';
              isValid = false;
          } else {
              document.getElementById('regPasswordError').textContent = '';
          }
          
          // 验证确认密码
          if (confirmPassword === '') {
              document.getElementById('regConfirmPasswordError').textContent = '请确认密码';
              isValid = false;
          } else if (password !== confirmPassword) {
              document.getElementById('regConfirmPasswordError').textContent = '两次输入的密码不一致';
              isValid = false;
          } else {
              document.getElementById('regConfirmPasswordError').textContent = '';
          }
          
          if (isValid) {
              alert('注册成功！');
              // 实际项目中这里会有AJAX请求到服务器保存用户数据
              window.location.href = 'login.html';
          }
      });
  }
  
  // 邮箱验证函数
  function validateEmail(email) {
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      return re.test(email);
  }
});

    // 检查用户是否已登录
    function checkLogin() {
        return localStorage.getItem('isLoggedIn') === 'true';
    }
    
    // 设置登录状态
    function setLoggedIn(isLoggedIn) {
        localStorage.setItem('isLoggedIn', isLoggedIn);
    }
    
    // 登录函数
    function login(username) {
        setLoggedIn(true);
        localStorage.setItem('username', username);
        updateUI();
    }
    
    // 登出函数
    function logout() {
        setLoggedIn(false);
        localStorage.removeItem('username');
        updateUI();
    }
    
    // 更新UI显示
    function updateUI() {
        const isLoggedIn = checkLogin();
        const username = localStorage.getItem('username');
        
        // 更新导航栏
        const navItems = document.querySelector('nav ul');
        if (isLoggedIn) {
            navItems.innerHTML = `
                <li><a href="index.html">Home</a></li>
                <li><a href="myblog.html">我的博客</a></li>
                <li><a href="#" onclick="logout()">退出登录</a></li>
            `;
        } else {
            navItems.innerHTML = `
                <li><a href="index.html">Home</a></li>
                <li><a href="login.html">登录</a></li>
                <li><a href="register.html">注册</a></li>
            `;
        }
    }
    // 修改登录逻辑
document.getElementById('login-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        if (response.ok) {
            // 保存令牌到localStorage
            localStorage.setItem('token', data.token);
            // 登录成功后获取用户信息
            fetchCurrentUser();
            window.location.href = 'myblog.html';
        } else {
            alert(data.error);
        }
    } catch (error) {
        console.error('登录失败:', error);
        alert('登录请求失败');
    }
});

// 获取当前用户信息
async function fetchCurrentUser() {
    const token = localStorage.getItem('token');
    if (!token) return;
    
    try {
        const response = await fetch('http://localhost:8080/api/user/me', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const user = await response.json();
        if (response.ok) {
            // 更新UI显示用户信息
            updateUserUI(user);
        } else {
            // 令牌无效，清除本地存储
            localStorage.removeItem('token');
        }
    } catch (error) {
        console.error('获取用户信息失败:', error);
    }
}

