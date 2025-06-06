import { Link } from 'react-router-dom';
import '../styles/Header.css';
import { jwtDecode } from 'jwt-decode';



export default function Header() {
    const token = localStorage.getItem('token');

    let role = null;
    if (token) {
        try {
            const decoded = jwtDecode(token);
            role = decoded.role;
        } catch (err) {
            console.error('Ошибка декодирования токена:', err);
        }
    }

    return (
        <header className="header">
            <nav className="nav">
                <Link to="/" className="logo">Чистый Дом</Link>
                <div className="nav-links">
                    <>
                        <Link to="/">О компании</Link>
                        <Link to="/services">Услуги</Link>
                        <Link to="/reviews">Отзывы</Link>
                    </>
                    {token ? (
                        <>
                            {role === 'admin' ? (
                                <Link to="/admin/orders">Заказы пользователей</Link>
                            ) : (
                                <Link to="/orders">Мои заявки</Link>
                            )}
                            <button onClick={() => {
                                localStorage.removeItem('token');
                                window.location.href = '/';
                            }}>Выйти</button>
                        </>
                    ) : (
                        <>
                            <Link to="/login">Вход</Link>
                            <Link to="/register">Регистрация</Link>
                        </>
                    )}
                </div>
            </nav>
        </header>
    );
}
