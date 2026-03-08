import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Header from './components/Layout/Header';
import SearchBar from './components/Home/SearchBar';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import Profile from './components/Profile/Profile';
import { useEffect, useState } from 'react';
import axios from 'axios';

function App() {
  const [tours, setTours] = useState([]);

  // Vẫn giữ phần gọi API Golang nhé
  useEffect(() => {
    axios.get('http://localhost:8080/api/tours')
      .then(res => setTours(res.data))
      .catch(err => console.error(err));
  }, []);

  return (
    <AuthProvider>
      <Router>
        <div className="app-container">
          <Header />
          
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="/" element={
              <main style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
                <SearchBar />
                
                <h2 style={{ marginTop: '40px' }}>Tour nổi bật</h2>
                <div style={{ display: 'flex', gap: '20px', marginTop: '20px' }}>
                  {tours.map(tour => (
                    <div key={tour.id} style={{
                      border: '1px solid #333',
                      borderRadius: '8px',
                      padding: '15px',
                      width: '300px',
                      backgroundColor: '#1e1e1e'
                    }}>
                      <div style={{ width: '100%', height: '150px', backgroundColor: '#333', borderRadius: '4px', marginBottom: '15px' }}></div>
                      <h3>{tour.name}</h3>
                      <p style={{ color: '#ff4d4f', fontWeight: 'bold' }}>{tour.price}</p>
                    </div>
                  ))}
                </div>
              </main>
            } />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;