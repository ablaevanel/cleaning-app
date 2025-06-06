import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api';
import Header from '../components/Header';
import Footer from '../components/Footer';
import '../styles/Register.css';

export default function Register() {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const [isError, setIsError] = useState(false);
    const navigate = useNavigate();

    const handleRegister = async () => {
        try {
            await api.post('/register', { name, email, password });
            setIsError(false);
            setMessage('Регистрация прошла успешно! Переход на страницу входа...');
            setTimeout(() => navigate('/login'), 2000);
        } catch (err) {
            setIsError(true);
            setMessage('Ошибка регистрации');
        }
    };

    return (
        <>
            <Header />
            <div className="register-container">
                <h2>Регистрация</h2>
                <input
                    type="text"
                    placeholder="Имя"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <input
                    type="email"
                    placeholder="Почта"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Пароль"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <button onClick={handleRegister}>Зарегистрироваться</button>
                {message && (
                    <p style={{ color: isError ? 'red' : 'green' }}>{message}</p>
                )}
            </div>
            <Footer />
        </>
    );
}
