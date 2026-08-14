import { LoginForm } from "@/components/forms/login-form";
import { Leaf } from "lucide-react";

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-green-50 to-white p-4">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center">
          <div className="flex justify-center mb-2">
            <div className="bg-green-100 p-3 rounded-full">
              <Leaf className="h-10 w-10 text-green-600" />
            </div>
          </div>
          <h1 className="text-3xl font-bold text-green-700">Gratia</h1>
          <p className="text-gray-600 mt-1">Food Management System</p>
        </div>
        <LoginForm />
      </div>
    </div>
  );
}