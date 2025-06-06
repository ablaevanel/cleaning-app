import { useEffect, useState } from 'react';
import api from '../api';
import Header from '../components/Header';
import Footer from '../components/Footer';
import '../styles/AdminOrders.css';

export default function AdminOrders() {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchOrders = async () => {
            try {
                const token = localStorage.getItem('token');
                const res = await api.get('/admin/orders', {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });
                setOrders(res.data);
            } catch (err) {
                console.error('Ошибка загрузки заказов:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchOrders();
    }, []);

    const handleStatusChange = async (orderId, newStatus) => {
        try {
            const token = localStorage.getItem('token');
            await api.patch(`/admin/orders/${orderId}`, { status: newStatus }, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            setOrders(prev =>
                prev.map(order =>
                    order.id === orderId ? { ...order, status: newStatus } : order
                )
            );
        } catch (err) {
            console.error('Ошибка обновления статуса:', err);
            alert('Не удалось обновить статус');
        }
    };

    const handleDelete = async (orderId) => {
        if (!window.confirm('Вы уверены, что хотите удалить этот заказ?')) return;

        try {
            const token = localStorage.getItem('token');
            await api.delete(`/admin/orders/${orderId}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            setOrders(prev => prev.filter(order => order.id !== orderId));
            alert('Заказ удалён');
        } catch (err) {
            console.error('Ошибка удаления заказа:', err);
            alert('Не удалось удалить заказ');
        }
    };

    return (
        <>
            <Header />
            <main className="admin-orders-container">
                <h2>Заказы пользователей</h2>
                {loading ? (
                    <p>Загрузка...</p>
                ) : orders.length === 0 ? (
                    <p>Заявки отсутствуют.</p>
                ) : (
                    <ul className="admin-orders-list">
                        {orders.map(order => (
                            <li key={order.id} className="admin-order-item">
                                <strong>Пользователь:</strong> {order.user}<br />
                                <strong>Услуга:</strong> {order.service}<br />
                                <strong>Статус:</strong>
                                <select
                                    value={order.status}
                                    onChange={(e) => handleStatusChange(order.id, e.target.value)}
                                >
                                    <option value="created">created</option>
                                    <option value="in progress">in progress</option>
                                    <option value="completed">completed</option>
                                </select>
                                <br />
                                <button
                                    className="delete-button"
                                    onClick={() => handleDelete(order.id)}
                                >
                                    Удалить
                                </button>
                            </li>
                        ))}
                    </ul>
                )}
            </main>
            <Footer />
        </>
    );
}
