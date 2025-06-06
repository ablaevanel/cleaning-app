import { useEffect, useState } from 'react';
import api from '../api';
import Header from '../components/Header';
import Footer from '../components/Footer';
import '../styles/Reviews.css';

export default function Reviews() {
    const [reviews, setReviews] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchReviews = async () => {
            try {
                const response = await api.get('/reviews');
                setReviews(response.data);
            } catch (err) {
                setError('Не удалось загрузить отзывы');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchReviews();
    }, []);

    const renderStars = (rating) => {
        return '★'.repeat(rating) + '☆'.repeat(5 - rating);
    };

    return (
        <>
            <Header />
            <main className="reviews-container">
                <h2>Все отзывы</h2>
                
                {loading ? (
                    <p>Загрузка отзывов...</p>
                ) : error ? (
                    <p className="error">{error}</p>
                ) : reviews.length === 0 ? (
                    <p>Отзывов пока нет</p>
                ) : (
                    <div className="reviews-grid">
                        {reviews.map(review => (
                            <div key={review.id} className="review-card">
                                <div className="review-header">
                                    <span className="review-service">{review.service_name}</span>
                                    <span className="review-rating">
                                        {renderStars(review.rating)}
                                    </span>
                                </div>
                                <div className="review-body">
                                    <p className="review-comment">{review.comment}</p>
                                </div>
                                <div className="review-footer">
                                    <span className="review-author">Автор: {review.user_name}</span>
                                    <span className="review-date">
                                        {new Date(review.created_at).toLocaleDateString()}
                                    </span>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </main>
            <Footer />
        </>
    );
}