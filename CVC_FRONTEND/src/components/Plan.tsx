import { useDispatch} from "react-redux";
import { useNavigate } from "react-router-dom";
import { logout } from "@/redux/Store/userSlice";

export const Plan = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    localStorage.removeItem("falcoRunning")
    dispatch(logout());
    navigate("/");
  };

  return (
    <div className='flex sticky top-[calc(100vh_-_48px_-_16px)] flex-col h-12 border-t px-2 border-stone-300 justify-end text-xs'>
      <div className='flex items-center justify-between'>
        <div>
          <p className='font-bold'>CVC 2025</p>
          <p className='text-stone-500'>Scanning Platform</p>
        </div>
        <button onClick={handleLogout} className='px-2 py-1.5 font-medium bg-stone-200 hover:bg-stone-300 transition-colors rounded'>
          Log Out
        </button>
      </div>
    </div>
  );
};
