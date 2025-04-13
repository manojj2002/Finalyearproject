import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import { useDispatch } from "react-redux";
import { setUser } from "@/redux/Store/userSlice";
import {jwtDecode} from "jwt-decode";

const GithubRedirectHandler = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();

  // useEffect(() => {
  //   const urlParams = new URLSearchParams(location.search);
  //   const token = urlParams.get("token");
  //   console.log(token)
  //   if (!token) {
  //     toast.error("❌ GitHub login failed: No token received.");
  //     navigate("/");
  //     return;
  //   }

  //   try {
  //     localStorage.setItem("token", token);
  //     const decoded: any = jwtDecode(token);
  //     const username = decoded.username || decoded.email || "user";
  //     localStorage.setItem("username", username);
  //     dispatch(setUser({ username: username, token }));
  //     toast.success("✅ GitHub login successful");
  //     navigate("/dashboard");
  //   } catch (err) {
  //     toast.error("❌ Failed to decode token.");
  //     navigate("/");
  //   }
  // }, [navigate, dispatch]);
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      const urlParams = new URLSearchParams(location.search);
      const token = urlParams.get("token");
  
      console.log("Token:", token);
      if (!token) {
        toast.error("❌ GitHub login failed: No token received.");
        navigate("/");
        return;
      }
  
      try {
        localStorage.setItem("token", token);
        const decoded: any = jwtDecode(token);
        const username = decoded.username || decoded.email || "user";
        localStorage.setItem("username", username);
        dispatch(setUser({ username: username, token }));
        toast.success("✅ GitHub login successful");
        navigate("/dashboard");
      } catch (err) {
        toast.error("❌ Failed to decode token.");
        navigate("/");
      }
    }, 1000); // Wait a second for the URL to settle
  
    return () => clearTimeout(timeoutId); // Clean up timeout
  }, [navigate, dispatch]);
  

  return <p className="text-center mt-10">🔄 Authenticating with GitHub...</p>;
};

export default GithubRedirectHandler;
