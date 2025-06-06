import { useState, useEffect } from 'react';
import api from '../api';
import Header from '../components/Header';
import Footer from '../components/Footer';
import '../styles/Orders.css';

export default function Orders() {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [activeReview, setActiveReview] = useState(null);
    const [reviewData, setReviewData] = useState({ rating: 5, comment: '' });
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchOrders = async () => {
            try {
                const res = await api.get('/orders');
                setOrders(res.data || []);
            } catch (err) {
                setError('Не удалось загрузить заказы');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };
        fetchOrders();
    }, []);

    const handleReviewSubmit = async (orderId) => {
        try {
            await api.post(`/orders/${orderId}/reviews`, reviewData);
            const updatedOrders = orders.map(order => 
                order.id === orderId ? { ...order, hasReview: true } : order
            );
            setOrders(updatedOrders);
            setActiveReview(null);
            setReviewData({ rating: 5, comment: '' });
        } catch (err) {
            setError(err.response?.data?.error || 'Ошибка при отправке отзыва');
        }
    };

    const handleDeleteOrder = async (orderId) => {
        if (!window.confirm('Вы уверены, что хотите удалить этот заказ?')) {
            return;
        }

        try {
            await api.delete(`/orders/${orderId}`);
            setOrders(orders.filter(order => order.id !== orderId));
        } catch (err) {
            setError(err.response?.data?.error || 'Не удалось удалить заказ');
        }
    };

    return (
        <>
            <Header />
            <main className="orders-container">
                <h2>Мои заявки</h2>
                {error && <p className="error">{error}</p>}
                
                {loading ? (
                    <p>Загрузка...</p>
                ) : orders.length === 0 ? (
                    <p>У вас пока нет заявок.</p>
                ) : (
                    <ul className="orders-list">
                        {orders.map(order => (
                            <li key={order.id} className="order-item">
                                <strong>Услуга:</strong> {order.service}<br />
                                <strong>Статус:</strong> {order.status}<br />
                                <strong>Дата:</strong> {new Date(order.created_at).toLocaleString()}

                                <div className="order-actions">
                                    {/* Кнопка удаления для заказов со статусом created */}
                                    {order.status === 'created' && (
                                        <button 
                                            onClick={() => handleDeleteOrder(order.id)}
                                            className="delete-btn"
                                        >
                                            Удалить заказ
                                        </button>
                                    )}

                                    {/* Форма отзыва для завершенных заказов */}
                                    {order.status === 'completed' && !order.hasReview && (
                                        <div className="review-section">
                                            {activeReview === order.id ? (
                                                <div className="review-form">
                                                    <select
                                                        value={reviewData.rating}
                                                        onChange={(e) => setReviewData({
                                                            ...reviewData,
                                                            rating: parseInt(e.target.value)
                                                        })}
                                                    >
                                                        {[1, 2, 3, 4, 5].map(num => (
                                                            <option key={num} value={num}>{num} ★</option>
                                                        ))}
                                                    </select>
                                                    <textarea
                                                        placeholder="Ваш отзыв..."
                                                        value={reviewData.comment}
                                                        onChange={(e) => setReviewData({
                                                            ...reviewData,
                                                            comment: e.target.value
                                                        })}
                                                        rows={3}
                                                    />
                                                    <div className="review-buttons">
                                                        <button 
                                                            onClick={() => handleReviewSubmit(order.id)}
                                                            className="submit-btn"
                                                        >
                                                            Отправить
                                                        </button>
                                                        <button 
                                                            onClick={() => setActiveReview(null)}
                                                            className="cancel-btn"
                                                        >
                                                            Отмена
                                                        </button>
                                                    </div>
                                                </div>
                                            ) : (
                                                <button 
                                                    onClick={() => setActiveReview(order.id)}
                                                    className="add-review-btn"
                                                >
                                                    Оставить отзыв
                                                </button>
                                            )}
                                        </div>
                                    )}
                                </div>
                            </li>
                        ))}
                    </ul>
                )}
            </main>
            <Footer />
        </>
    );
}