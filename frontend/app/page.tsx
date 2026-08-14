"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Leaf, Users, Heart, Clock } from "lucide-react";
import { motion } from "framer-motion";

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-green-50 to-white">
      {/* Navigation */}
      <nav className="border-b bg-white/50 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <Leaf className="h-8 w-8 text-green-600" />
            <span className="text-2xl font-bold text-green-700">Gratia</span>
          </div>
          <div className="flex items-center gap-4">
            <Link href="/login">
              <Button variant="outline">Log In</Button>
            </Link>
            <Link href="/register">
              <Button>Get Started</Button>
            </Link>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
            Bridge the Gap Between
            <br />
            <span className="text-green-600">Excess Food</span> and
            <span className="text-green-600"> Those in Need</span>
          </h1>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto mb-8">
            Gratia connects restaurants with excess food to NGOs and communities,
            reducing waste and fighting hunger together.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link href="/register">
              <Button size="lg" className="text-lg px-8">
                Start Helping
              </Button>
            </Link>
            <Link href="#how-it-works">
              <Button size="lg" variant="outline" className="text-lg px-8">
                Learn More
              </Button>
            </Link>
          </div>
        </motion.div>

        {/* Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mt-20">
          <div className="bg-white p-6 rounded-xl shadow-sm hover:shadow-md transition">
            <div className="text-3xl font-bold text-green-600">100+</div>
            <div className="text-gray-600">Restaurants</div>
          </div>
          <div className="bg-white p-6 rounded-xl shadow-sm hover:shadow-md transition">
            <div className="text-3xl font-bold text-green-600">50+</div>
            <div className="text-gray-600">NGOs</div>
          </div>
          <div className="bg-white p-6 rounded-xl shadow-sm hover:shadow-md transition">
            <div className="text-3xl font-bold text-green-600">10K+</div>
            <div className="text-gray-600">Meals Served</div>
          </div>
          <div className="bg-white p-6 rounded-xl shadow-sm hover:shadow-md transition">
            <div className="text-3xl font-bold text-green-600">5K+</div>
            <div className="text-gray-600">People Helped</div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section id="how-it-works" className="container mx-auto px-4 py-20">
        <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">
          How Gratia Works
        </h2>
        <div className="grid md:grid-cols-3 gap-8">
          <motion.div
            whileHover={{ y: -5 }}
            className="text-center p-6 bg-white rounded-xl shadow-sm"
          >
            <div className="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <Users className="h-8 w-8 text-green-600" />
            </div>
            <h3 className="text-xl font-semibold mb-2">1. Register</h3>
            <p className="text-gray-600">Join as a donor, NGO, or admin</p>
          </motion.div>
          <motion.div
            whileHover={{ y: -5 }}
            className="text-center p-6 bg-white rounded-xl shadow-sm"
          >
            <div className="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <Clock className="h-8 w-8 text-green-600" />
            </div>
            <h3 className="text-xl font-semibold mb-2">2. List Food</h3>
            <p className="text-gray-600">Donors list excess food items</p>
          </motion.div>
          <motion.div
            whileHover={{ y: -5 }}
            className="text-center p-6 bg-white rounded-xl shadow-sm"
          >
            <div className="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <Heart className="h-8 w-8 text-green-600" />
            </div>
            <h3 className="text-xl font-semibold mb-2">3. Claim & Deliver</h3>
            <p className="text-gray-600">NGOs claim and distribute food</p>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t bg-white">
        <div className="container mx-auto px-4 py-8 text-center text-gray-600">
          <p>© {new Date().getFullYear()} Gratia. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}