import { useState, useEffect } from 'react';

const HEALTH_CHECK_INTERVAL = 5 * 60 * 1000; // 5分鐘

const SERVICES = [
    {
        name: 'Auth Service',
        url: 'https://ordermanagersystem-auth-service.onrender.com/health'
    },
    {
        name: 'Product Service',
        url: 'https://ordermanagersystem-product-service.onrender.com/health'
    },
    {
        name: 'Payment Service',
        url: 'https://ordermanagersystem-payment-service.onrender.com/health'
    },
    {
        name: 'Notification Service',
        url: 'https://ordermanagersystem-notification-service.onrender.com/health'
    },
    {
        name: 'Cart Service',
        url: 'https://ordermanagersystem.onrender.com/health'
    }
];

export const useHealthCheck = () => {
    const [serviceStatus, setServiceStatus] = useState({});
    const [isAnyServiceDown, setIsAnyServiceDown] = useState(false);

    const checkHealth = async () => {
        const newStatus = {};
        let anyServiceDown = false;

        await Promise.all(
            SERVICES.map(async (service) => {
                try {
                    const response = await fetch(service.url);
                    const isHealthy = response.ok;
                    newStatus[service.name] = isHealthy;
                    if (!isHealthy) anyServiceDown = true;
                } catch (error) {
                    console.error(`Health check failed for ${service.name}:`, error);
                    newStatus[service.name] = false;
                    anyServiceDown = true;
                }
            })
        );

        setServiceStatus(newStatus);
        setIsAnyServiceDown(anyServiceDown);
    };

    useEffect(() => {
        // 初始檢查
        checkHealth();

        // 設置定期檢查
        const interval = setInterval(checkHealth, HEALTH_CHECK_INTERVAL);

        // 清理函數
        return () => clearInterval(interval);
    }, []);

    return {
        serviceStatus,
        isAnyServiceDown
    };
}; 