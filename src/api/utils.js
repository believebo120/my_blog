import axios from 'axios';

// 文件上传
export const uploadFile = (file) => {
    const formData = new FormData();
    formData.append('file', file);
    return axios.post('/uploads', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};