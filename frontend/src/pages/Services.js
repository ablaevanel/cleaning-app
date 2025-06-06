import { useEffect, useState } from 'react';
import api from '../api';
import Header from '../components/Header';
import Footer from '../components/Footer';
import '../styles/Services.css';

export default function Services() {
    const [services, setServices] = useState([]);

    useEffect(() => {
        const fetchServices = async () => {
            try {
                const res = await api.get('/services');
                setServices(res.data);
            } catch (err) {
                console.error('Ошибка загрузки услуг:', err);
            }
        };
        fetchServices();
    }, []);

    const handleOrder = async (serviceId) => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                alert('Сначала войдите в аккаунт');
                return;
            }

            const response = await api.post(
                '/orders',
                { service_id: serviceId },
                {
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${token}`
                    }
                }
            );

            if (response.data && response.data.order_id) {
                alert(`Заявка #${response.data.order_id} успешно оформлена`);
            } else {
                alert('Заявка оформлена, но не получен ID заказа');
            }
        } catch (err) {
            console.error('Ошибка оформления заявки:', err);
            alert(err.response?.data?.error || 'Не удалось оформить заявку');
        }
    };

    return (
        <>
            <Header />
            <main className="services-container">
                <h2>Наши услуги</h2>
                <div className="services-list">
                    {services.map(service => (
                        <div key={service.id} className="service-card">
                            <img src={service.image_url} alt={service.name} className="service-image" />
                            <h3>{service.name}</h3>
                            <p>{service.description}</p>
                            <p className="service-price">{service.price} ₽</p>
                            <button
                                className="order-button"
                                onClick={() => handleOrder(service.id)}
                            >
                                Оформить заявку
                            </button>
                        </div>
                    ))}
                </div>
            </main>
            <Footer />
        </>
    );
}
