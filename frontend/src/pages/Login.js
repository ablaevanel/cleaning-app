import { useState } from 'react';
import axios from '../api';
import '../styles/Login.css';
import Header from '../components/Header';
import Footer from '../components/Footer';

export default function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const res = await axios.post('/login', { email, password });
            const token = res.data.token;
            localStorage.setItem('token', token);
            window.location.href = '/';
        } catch (err) {
            setError('Неверный email или пароль');
        }
    };

    return (
        <>
            <Header />
            <div className="login-container">
                <h2>Вход</h2>
                <form onSubmit={handleLogin}>
                    <input
                        type="email"
                        placeholder="Email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                    <button type="submit">Войти</button>
                </form>
                {error && <p className="error-message">{error}</p>}
            </div>
            <Footer />
        </>
    );
}