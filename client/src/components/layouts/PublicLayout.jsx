import { Outlet } from 'react-router-dom';
import Navbar from '../common/Navbar';
import Footer from '../common/Footer';

const PublicLayout = () => {
  return (
    <div className='min-h-screen flex flex-col'>
      <div className='sticky top-0 z-40'>
        <Navbar />
      </div>

      <main className='flex-1'>
        <Outlet />
      </main>

      <Footer />
    </div>
  );
};

export default PublicLayout;
