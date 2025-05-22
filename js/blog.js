document.addEventListener('DOMContentLoaded', function() {
    // 搜索功能
    const searchBtn = document.getElementById('searchBtn');
    if (searchBtn) {
        searchBtn.addEventListener('click', function() {
            const searchTerm = document.getElementById('searchInput').value.trim();
            if (searchTerm) {
                alert(`搜索: ${searchTerm}`);
            }
        });
    }
    
    // 删除文章功能
    const deleteButtons = document.querySelectorAll('.btn-danger');
    deleteButtons.forEach(button => {
        button.addEventListener('click', function() {
            if (confirm('确定要删除这篇文章吗？')) {
                const post = this.closest('.post');
                post.style.opacity = '0';
                setTimeout(() => {
                    post.remove();
                    alert('文章已删除');
                }, 300);
            }
        });
    });
    
    // 添加文章表单
    const addPostForm = document.getElementById('addPostForm');
    if (addPostForm) {
        addPostForm.addEventListener('submit', function(e) {
            e.preventDefault();
            alert('文章已添加');
            window.location.href = 'blog.html';
        });
    }
    
    // 编辑文章表单
    const editPostForm = document.getElementById('editPostForm');
    if (editPostForm) {
        editPostForm.addEventListener('submit', function(e) {
            e.preventDefault();
            alert('文章已更新');
            window.location.href = 'blog.html';
        });
    }
});