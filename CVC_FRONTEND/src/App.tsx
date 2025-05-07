import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { BrowserRouter, Routes, Route } from 'react-router-dom';


import { Home } from './pages/Home';
import { DashboardHome } from './pages/DashboardHome';
import { DashboardLayout } from './pages/DashboardLayout';
import { ScanResultsPage } from "./pages/ScanResultsPage";
import { SignupForm } from "./pages/UserRegister";
import GithubRedirectHandler from "./pages/GithubRedirectHandler";
import { LoginForm } from "./pages/LoginPage";
const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path='/' element={<Home />} />

        <Route path='/register' element={<SignupForm/>} />
        <Route path= '/login' element={<LoginForm/>}/>
        <Route path="/githubLogin" element={<GithubRedirectHandler />} />

        <Route path='/dashboard' element={<DashboardLayout />}>
          <Route index element={<DashboardHome />} />
          <Route path='result' element={<ScanResultsPage />} />
        </Route>
      </Routes>

      <ToastContainer position="top-right" autoClose={3000} />
    </BrowserRouter>
  );
};

export default App;
