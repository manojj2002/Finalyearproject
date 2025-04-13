import { useState } from "react";
import { FaGithub } from "react-icons/fa";
import { toast } from "react-toastify";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useDispatch} from "react-redux";
import { setUser } from "@/redux/Store/userSlice";



const RegisterOrGithub = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch()
    const [name, setName] = useState("")
    const [email, setEmail] = useState("");
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await axios.post("/api/user/register", {
        name,
        email,
        username,
        password,
      });
      dispatch(setUser({
        token: res.data.token,
        username: res.data.username,
      }));
      console.log(username)
      toast.success("🎉 Registered successfully!");
      navigate("/dashboard");
      setName("")
      setEmail("");
      setUsername("");
      setPassword("");
    } catch (err) {
      toast.error("❌ Registration failed");
    }
  };

  const handleGithubLogin = () => {
    window.location.href = "/api/auth/login";
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4">
      <div className="bg-white shadow-md rounded-lg p-8 max-w-md w-full">
        <h2 className="text-2xl font-bold mb-6 text-center text-black">Create Account</h2>
        
        <form onSubmit={handleRegister} className="space-y-4">
        <input
            type="text"
            placeholder="Name"
            className="w-full p-2 border rounded bg-blue-100 text-black"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
          <input
            type="email"
            placeholder="Email"
            className="w-full p-2 border rounded bg-blue-100 text-black"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="text"
            placeholder="Username"
            className="w-full p-2 border rounded bg-blue-100 text-black"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Password"
            className="w-full p-2 border rounded bg-blue-100 text-black"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button
           onClick={handleRegister}
            type="submit"
            className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
          >
            Register
          </button>
        </form>

        <div className="my-4 text-center text-gray-500">or</div>

        <button
          onClick={handleGithubLogin}
          className="w-full flex items-center justify-center bg-black text-white py-2 rounded hover:bg-gray-800"
        >
          <FaGithub className="mr-2" /> Sign in with GitHub
        </button>
      </div>
    </div>
  );
};

export default RegisterOrGithub;
