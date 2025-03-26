import axios from 'axios';

const AUTH_SERVICE_URL = 'https://ordermanagersystem-auth-service.onrender.com';

// 刷新 token 的函數
export const refreshToken = async () => {
    try {
        const refreshToken = localStorage.getItem('refreshToken');
        if (!refreshToken) {
            throw new Error('沒有可用的刷新令牌');
        }

        const response = await axios.post(`${AUTH_SERVICE_URL}/api/v1/auth/refresh`, {
            refreshToken: refreshToken
        });

        if (response.data.token) {
            localStorage.setItem('userToken', response.data.token);
            if (response.data.refreshToken) {
                localStorage.setItem('refreshToken', response.data.refreshToken);
            }
            return response.data.token;
        } else {
            throw new Error('刷新令牌失敗');
        }
    } catch (error) {
        console.error('刷新令牌失敗:', error);
        localStorage.removeItem('userToken');
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('userData');
        throw error;
    }
};

// 建立一個帶有統一錯誤處理的 axios 實例
export const createAuthAxios = (navigate) => {
    const instance = axios.create({
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        }
    });

    // 是否正在刷新 token
    let isRefreshing = false;
    // 等待令牌刷新的請求隊列
    let failedQueue = [];

    // 處理隊列中的請求
    const processQueue = (error, token = null) => {
        failedQueue.forEach(prom => {
            if (error) {
                prom.reject(error);
            } else {
                prom.resolve(token);
            }
        });
        failedQueue = [];
    };

    // 請求攔截器 - 添加令牌
    instance.interceptors.request.use(
        (config) => {
            const token = localStorage.getItem('userToken');
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
            return config;
        },
        (error) => Promise.reject(error)
    );

    // 響應攔截器 - 處理令牌刷新
    instance.interceptors.response.use(
        (response) => response,
        async (error) => {
            const originalRequest = error.config;

            // 如果是 401 錯誤且不是刷新令牌的請求，且請求未重試過
            if (error.response?.status === 401
                && !originalRequest._retry
                && !originalRequest.url.includes('/api/v1/auth/refresh')) {

                if (isRefreshing) {
                    // 如果已經在刷新，將請求加入隊列
                    return new Promise((resolve, reject) => {
                        failedQueue.push({ resolve, reject });
                    })
                        .then(token => {
                            originalRequest.headers.Authorization = `Bearer ${token}`;
                            return instance(originalRequest);
                        })
                        .catch(err => Promise.reject(err));
                }

                originalRequest._retry = true;
                isRefreshing = true;

                try {
                    // 嘗試刷新令牌
                    const newToken = await refreshToken();
                    originalRequest.headers.Authorization = `Bearer ${newToken}`;

                    // 處理等待中的請求
                    processQueue(null, newToken);
                    return instance(originalRequest);
                } catch (refreshError) {
                    // 刷新失敗，處理等待中的請求
                    processQueue(refreshError);

                    // 重定向到登入頁面
                    if (navigate) {
                        navigate('/login');
                    }
                    return Promise.reject(refreshError);
                } finally {
                    isRefreshing = false;
                }
            }

            return Promise.reject(error);
        }
    );

    return instance;
}; 