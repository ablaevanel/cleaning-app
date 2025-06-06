import Header from '../components/Header';
import Footer from '../components/Footer';
import '../styles/Home.css';

export default function Home() {
    return (
        <div className="home-wrapper">
            <Header />
            <main className="home-container">
                <div className="home-content">
                    <div className="home-text">
                        <h1>О компании</h1>
                        <p>
                            Мы — профессиональная клининговая компания «Чистый Дом». Предлагаем качественные услуги уборки для квартир, домов и офисов.
                        </p>
                        <p>
                            В нашей команде — опытные специалисты, которые используют современное оборудование и безопасные моющие средства.
                        </p>
                        <p>
                            Чистота, пунктуальность и доверие — вот наши приоритеты.
                        </p>
                    </div>
                    <div className="home-image">
                        <img src="https://super-cleaning-spb.ru/uploads/s/g/f/l/gflxywmvkvrh/img/full_gpJHMeae.jpg" alt="Клининг" />
                    </div>
                </div>
            </main>
            <Footer />
        </div>
    );
}
