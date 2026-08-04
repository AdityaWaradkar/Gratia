import Link from 'next/link'
import { Leaf } from 'lucide-react'

export function Footer() {
  return (
    <footer className="border-t bg-white">
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <div className="flex items-center gap-2">
            <Leaf className="h-5 w-5 text-green-600" />
            <span className="font-semibold text-green-700">Gratia</span>
            <span className="text-sm text-gray-500">© 2026</span>
          </div>
          <div className="flex gap-6 text-sm text-gray-600">
            <Link href="/about" className="hover:text-green-600 transition">
              About
            </Link>
            <Link href="/contact" className="hover:text-green-600 transition">
              Contact
            </Link>
            <Link href="/privacy" className="hover:text-green-600 transition">
              Privacy
            </Link>
            <Link href="/terms" className="hover:text-green-600 transition">
              Terms
            </Link>
          </div>
        </div>
      </div>
    </footer>
  )
}