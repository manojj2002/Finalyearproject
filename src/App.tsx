
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Home } from './pages/Home';
import { DashboardHome } from './pages/DashboardHome';
import { DashboardLayout } from './pages/DashboardLayout';

const App = () => {
  return (
   

    <BrowserRouter>
      <Routes>
        <Route
          path='/'
          element={<Home />}
        />
        <Route
          path='/dashboard'
          element={<DashboardLayout />}
        >
          <Route
            index
            element={<DashboardHome />}
          />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
