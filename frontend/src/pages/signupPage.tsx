import { useState } from "react";
import type { FormEvent } from "react";
import AuthImagePattern from "../components/authImagePattern";
import { userAuthService } from "../services/service_userAuth";
import type { SignUpPageDataContextType } from "../configs/types";

import toast from "react-hot-toast";
import { Link } from "react-router-dom";
import { Eye, EyeOff, Loader2, Lock, Mail, MessageSquare, User} from "lucide-react";

// Import the phone number input component and its CSS
import PhoneInput, { isValidPhoneNumber,  parsePhoneNumber } from "react-phone-number-input";
import "react-phone-number-input/style.css";

const signUpPageDataInitial: SignUpPageDataContextType = {
    first_name: "",
    last_name: "",
    email: "",
    password: "",
    phone: "",
}

const SignUpPage: React.FC = () => {
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const [formData, setFormData] = useState<SignUpPageDataContextType>({...signUpPageDataInitial});
    const { signUp, isSigningUp } = userAuthService();

    const [customFormattedPhoneNumber, setCustomFormattedPhoneNumber] = useState<string>("");

        // New: Custom handler for phone number formatting
    // Handler for the PhoneInput component
    const handlePhoneNumberChange = (value?: string) => {
        // Store the full international number from PhoneInput in the form state
        setFormData(prevData => ({ ...prevData, phone: value?.replace(/\s/g, '') || '' }));

        // Now, apply your custom formatting to the number for display purposes
        if (value) {
            try {
                const phoneNumber = parsePhoneNumber(value);
                if (phoneNumber) {
                    const nationalNumber = phoneNumber.nationalNumber;
                    const rawValue = nationalNumber.replace(/\D/g, '');
                    const truncatedValue = rawValue.substring(0, 9);
                    let formattedValue = '';
                    if (truncatedValue.length > 6) {
                        formattedValue = `${truncatedValue.slice(0, 3)} ${truncatedValue.slice(3, 6)} ${truncatedValue.slice(6, 9)}`;
                    } else if (truncatedValue.length > 3) {
                        formattedValue = `${truncatedValue.slice(0, 3)} ${truncatedValue.slice(3, 6)}`;
                    } else {
                        formattedValue = truncatedValue;
                    }
                    setCustomFormattedPhoneNumber(formattedValue);
                } else {
                    setCustomFormattedPhoneNumber("");
                }
            } catch (error) {
                console.error("Error formatting phone number:", error);
                setCustomFormattedPhoneNumber("");
            }
        } else {
            setCustomFormattedPhoneNumber("");
        }
    };
    const validateForm = () : boolean => {
        if (!formData.first_name || !formData.last_name || !formData.email || !formData.password) return false;
        if (!formData.first_name.trim() || !formData.last_name.trim()) {
            toast.error("First name and last name is required");
            return false;
        }
        if (!formData.email.trim()) {
            toast.error("Email is required");
            return false;
        }
        // Basic email format validation using a regular expression
        if (!/\S+@\S+\.\S+/.test(formData.email)) {
            toast.error("Invalid email format");
            return false;
        }

        // Updated Validation: Check for phone number
        if (!formData.phone) {
            toast.error("Phone number is required");
            return false;
        }
        // Use the library's validator for a more robust check
        if (!isValidPhoneNumber(formData.phone)) {
            toast.error("Invalid phone number");
            return false;
        }

        // Updated: Validate the phone number format
        // Remove spaces to check the raw length
        const rawPhoneNumber = formData.phone.replace(/\s/g, '');
        console.log("Raw phone number: ", rawPhoneNumber);
        if (rawPhoneNumber.length !== 12) {
            toast.error("Phone number must be exactly 9 digits.");
            return false;
        }
        // Optional: Validate length of national number for Vietnam
        const parsed = parsePhoneNumber(formData.phone);
        console.log("Parsed phone number: ", parsed?.country, parsed?.nationalNumber.length);
        if (parsed?.country === "VN" && parsed?.nationalNumber.length !== 9) {
            toast.error("Vietnam phone number must be exactly 9 digits");
            return false;
        }

        if (!formData.password) {
            toast.error("Password is required");
            return false;
        }
        if (formData.password.length < 6) {
            toast.error("Password must be at least 6 characters");
            return false;
        }

        return true;
    };

    const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!validateForm()) return;
        signUp(formData);
    };

    return (
        <div className="min-h-screen grid lg:grid-cols-2">
        {/* left side */}
        <div className="flex flex-col justify-center items-center p-6 sm:p-12">
            <div className="w-full max-w-md space-y-8">
            {/* LOGO */}
            <div className="text-center mb-8">
                <div className="flex flex-col items-center gap-2 group">
                <div
                    className="size-12 rounded-xl bg-primary/10 flex items-center justify-center 
                group-hover:bg-primary/20 transition-colors"
                >
                    <MessageSquare className="size-6 text-primary" />
                </div>
                <h1 className="text-2xl font-bold mt-2">Create Account</h1>
                <p className="text-base-content/60">Get started with your free account</p>
                </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
                <div className="form-control">
                <label className="label">
                    <span className="label-text font-medium">First Name</span>
                </label>
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <User className="size-5 text-base-content/40" />
                    </div>
                    <input
                    type="text"
                    className={`input input-bordered w-full pl-10`}
                    placeholder="John"
                    value={formData.first_name}
                    onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                    />
                </div>
                </div>

                <div className="form-control">
                <label className="label">
                    <span className="label-text font-medium">Last Name</span>
                </label>
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <User className="size-5 text-base-content/40" />
                    </div>
                    <input
                    type="text"
                    className={`input input-bordered w-full pl-10`}
                    placeholder="Doe"
                    value={formData.last_name}
                    onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                    />
                </div>
                </div>


                {/* New Phone Number Input Field */}
                <div className="form-control">
                    <label className="label">
                        <span className="label-text font-medium">Phone Number</span>
                    </label>
                    <div className="relative">
                        {/* The PhoneInput component */}
                        <PhoneInput
                            international
                            displayInitialValueAsLocalNumber = {false}
                            limitMaxLength = {true}
                            maxLength={14}
                            defaultCountry="VN" // Optional: Set a default country
                            placeholder="Enter phone number - 123 456 789"
                            value={formData.phone}
                            onChange={handlePhoneNumberChange}
                            className="input-phone-number" // Custom class for styling
                            />
                        {customFormattedPhoneNumber && (
                            <p className="text-sm text-base-content/60 mt-2">
                                Formatted: {customFormattedPhoneNumber}
                            </p>
                        )}
                    </div>
                </div>


                <div className="form-control">
                <label className="label">
                    <span className="label-text font-medium">Email</span>
                </label>
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Mail className="size-5 text-base-content/40" />
                    </div>
                    <input
                    type="email"
                    className={`input input-bordered w-full pl-10`}
                    placeholder="you@example.com"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    />
                </div>
                </div>

                <div className="form-control">
                <label className="label">
                    <span className="label-text font-medium">Password</span>
                </label>
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Lock className="size-5 text-base-content/40" />
                    </div>
                    <input
                    type={showPassword ? "text" : "password"}
                    className={`input input-bordered w-full pl-10`}
                    placeholder="••••••••"
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                    />
                    <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                    onClick={() => setShowPassword(!showPassword)}
                    >
                    {showPassword ? (
                        <EyeOff className="size-5 text-base-content/40" />
                    ) : (
                        <Eye className="size-5 text-base-content/40" />
                    )}
                    </button>
                </div>
                </div>

                <button type="submit" className="btn btn-primary w-full" disabled={isSigningUp}>
                {isSigningUp ? (
                    <>
                    <Loader2 className="size-5 animate-spin" />
                    Loading...
                    </>
                ) : (
                    "Create Account"
                )}
                </button>
            </form>

            <div className="text-center">
                <p className="text-base-content/60">
                Already have an account?{" "}
                <Link to="/login" className="link link-primary">
                    Sign in
                </Link>
                </p>
            </div>
            </div>
        </div>

        {/* right side */}

        <AuthImagePattern
            title="Join our community"
            subtitle="Connect with friends, share moments, and stay in touch with your loved ones."
        />
        </div>
    );
};

export default SignUpPage;