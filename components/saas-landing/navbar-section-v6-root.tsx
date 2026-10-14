"use client"

import { useRef, useState } from "react"

import { ArrowUpRightIcon, ChevronDownIcon, XIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet"
import { cn } from "@/lib/utils"

const loginHref = "/login"

type SolutionsItem = {
  label: string
  href: string
}

const SOLUTIONS_ITEMS: SolutionsItem[] = [
  { label: "Alur Kerja", href: "#" },
  { label: "Integrasi", href: "#" },
  { label: "Otomasi", href: "#" },
  { label: "Enterprise", href: "#" },
]

const NAV_LINKS = [
  { label: "Agen", href: "#" },
  { label: "Harga", href: "#" },
  { label: "Blog", href: "#" },
]

const DOCS_HREF = "/docs"

const dropdownSurfaceClass =
  "animate-dropdown-in absolute top-full z-60 mt-2 max-h-[min(70vh,24rem)] rounded-xl border border-border bg-background/98 shadow-md shadow-black/4 backdrop-blur-md dark:shadow-black/20"

const dropdownItemClass =
  "block whitespace-nowrap rounded-lg px-2.5 py-1.5 font-mono text-sm uppercase leading-none text-foreground/90 transition-colors hover:bg-muted hover:text-foreground"

const sheetLinkClass =
  "block py-1.5 font-mono text-[1.6875rem] font-medium uppercase tracking-tight text-zinc-50 leading-snug transition-colors hover:text-zinc-400"

const sheetSubLinkClass =
  "block rounded-md py-1 font-mono text-base uppercase leading-snug text-zinc-400 transition-colors hover:text-zinc-50"

function HamburgerIcon({ open }: { open: boolean }) {
  return (
    <div className="flex size-4 flex-col items-center justify-center gap-1">
      <span
        className={cn(
          "block h-px w-4 bg-current transition-all duration-300 ease-in-out",
          open ? "translate-y-[5px] rotate-45" : ""
        )}
      />
      <span
        className={cn(
          "block h-px w-4 bg-current transition-all duration-300 ease-in-out",
          open ? "scale-x-0 opacity-0" : ""
        )}
      />
      <span
        className={cn(
          "block h-px w-4 bg-current transition-all duration-300 ease-in-out",
          open ? "-translate-y-[5px] -rotate-45" : ""
        )}
      />
    </div>
  )
}

const linkClass =
  "font-mono text-sm uppercase text-muted-foreground transition-colors hover:text-foreground"

export function NavbarSectionV6Root({ logoSrc }: { logoSrc: string }) {
  const [menuOpen, setMenuOpen] = useState(false)
  const [solutionsOpen, setSolutionsOpen] = useState(false)
  const hoverPanelTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(
    undefined
  )

  const cancelHoverClose = () => {
    clearTimeout(hoverPanelTimerRef.current)
  }

  const openSolutions = () => {
    cancelHoverClose()
    setSolutionsOpen(true)
  }

  const closeSolutions = () => {
    cancelHoverClose()
    hoverPanelTimerRef.current = setTimeout(() => setSolutionsOpen(false), 120)
  }

  const closeMenu = () => setMenuOpen(false)

  return (
    <header className="sticky top-0 z-40 border-b border-transparent bg-background/80 backdrop-blur transition-colors duration-300">
      <style>{`
        @keyframes dropdown-in {
          from { opacity: 0; transform: translateY(-6px); }
          to   { opacity: 1; transform: translateY(0); }
        }
        .animate-dropdown-in {
          animation: dropdown-in 0.18s cubic-bezier(0.4,0,0.2,1) both;
        }
      `}</style>

      <Sheet open={menuOpen} onOpenChange={setMenuOpen}>
        <SheetContent
          aria-describedby={undefined}
          showCloseButton={false}
          side="right"
          className={cn(
            "z-50 flex h-dvh max-h-dvh min-h-dvh flex-col overflow-y-auto overscroll-none rounded-none border-0 p-0 shadow-none sm:rounded-none",
            "data-[side=right]:inset-x-0! data-[side=right]:inset-y-0! data-[side=right]:top-0! data-[side=right]:right-0! data-[side=right]:bottom-0! data-[side=right]:left-0!",
            "data-[side=right]:w-screen! data-[side=right]:max-w-none! data-[side=right]:sm:max-w-none!",
            "border-l-0!",
            "bg-zinc-950 text-zinc-50"
          )}
          style={{
            width: "100vw",
            maxWidth: "100vw",
            left: 0,
            right: 0,
            top: 0,
            bottom: 0,
          }}
          onInteractOutside={(e) => e.preventDefault()}
          onPointerDownOutside={(e) => e.preventDefault()}
        >
          <SheetHeader className="sr-only shrink-0 p-6">
            <SheetTitle>Navigasi</SheetTitle>
          </SheetHeader>

          <SheetClose asChild>
            <button
              type="button"
              className="absolute top-8 right-8 z-10 inline-flex rounded-md p-2 text-zinc-50 opacity-70 transition-opacity hover:opacity-100 focus-visible:ring-2 focus-visible:ring-zinc-400 focus-visible:ring-offset-2 focus-visible:ring-offset-zinc-950 focus-visible:outline-none"
              aria-label="Tutup menu"
            >
              <XIcon className="size-6" strokeWidth={1.5} aria-hidden />
            </button>
          </SheetClose>

          <nav
            aria-label="Utama"
            className="flex flex-1 flex-col justify-start px-8 pt-[24dvh]"
          >
            <ul className="flex flex-col gap-2">
              <li className="pt-1">
                <p className={sheetLinkClass}>Solusi</p>
                <ul className="mt-3 flex flex-col gap-2">
                  {SOLUTIONS_ITEMS.map(({ label, href }) => (
                    <li key={label}>
                      <a
                        className={sheetSubLinkClass}
                        href={href}
                        onClick={(e) => {
                          e.preventDefault()
                          closeMenu()
                        }}
                      >
                        {label}
                      </a>
                    </li>
                  ))}
                </ul>
              </li>

              {NAV_LINKS.map(({ label, href }) => (
                <li key={label}>
                  <a
                    className={sheetLinkClass}
                    href={href}
                    onClick={(e) => {
                      e.preventDefault()
                      closeMenu()
                    }}
                  >
                    {label}
                  </a>
                </li>
              ))}

              <li>
                <a
                  className={cn(
                    sheetLinkClass,
                    "inline-flex items-center gap-2"
                  )}
                  href={DOCS_HREF}
                  rel="noopener noreferrer"
                  target="_blank"
                  onClick={() => closeMenu()}
                >
                  Dokumentasi
                  <ArrowUpRightIcon aria-hidden className="size-5 shrink-0" />
                  <span className="sr-only"> (terbuka di tab baru)</span>
                </a>
              </li>
            </ul>
          </nav>

          <div className="mt-auto flex flex-col gap-3 px-8 pt-6 pb-[max(2rem,env(safe-area-inset-bottom,0px))]">
            <Button
              asChild
              variant="outline"
              className="h-11 w-full border-zinc-700 bg-transparent text-zinc-50 hover:bg-zinc-900 hover:text-zinc-50"
            >
              <a
                href={loginHref}
                className="font-mono uppercase"
                onClick={() => closeMenu()}
              >
                Masuk
              </a>
            </Button>
            <Button asChild className="h-11 w-full">
              <a
                href={loginHref}
                className="font-mono uppercase"
                onClick={() => closeMenu()}
              >
                Mulai Gratis
              </a>
            </Button>
          </div>
        </SheetContent>
      </Sheet>

      <div className="relative mx-auto flex h-16 w-full max-w-6xl items-center justify-between px-6">
        <div className="flex min-w-0 flex-1 items-center">
          <a
            aria-label="Beranda CipinCode CLI"
            className="inline-flex min-w-0 items-center gap-2 font-medium tracking-tight"
            href="/"
          >
            <img
              alt="CipinCode Logo"
              src={logoSrc}
              className="h-7 w-auto shrink-0 object-contain"
            />
          </a>
        </div>

        <nav className="absolute top-1/2 left-1/2 z-10 hidden min-w-0 -translate-x-1/2 -translate-y-1/2 items-center gap-6 md:flex">
          <div
            className="relative w-fit shrink-0"
            onMouseEnter={openSolutions}
            onMouseLeave={closeSolutions}
          >
            <button
              type="button"
              aria-expanded={solutionsOpen}
              aria-haspopup="menu"
              className={cn(
                linkClass,
                "inline-flex cursor-default items-center gap-1 rounded-md"
              )}
              onClick={(e) => e.preventDefault()}
            >
              Solutions
              <ChevronDownIcon aria-hidden className="size-3.5 shrink-0" />
            </button>

            {solutionsOpen && (
              <div
                className={cn(
                  dropdownSurfaceClass,
                  "left-0 max-w-[calc(100vw-2rem)] min-w-38 overflow-auto overflow-x-hidden p-0.5"
                )}
                onMouseEnter={openSolutions}
                onMouseLeave={closeSolutions}
              >
                <div className="flex flex-col py-0.5">
                  {SOLUTIONS_ITEMS.map((item) => (
                    <a
                      key={item.label}
                      className={dropdownItemClass}
                      href={item.href}
                      onClick={(e) => e.preventDefault()}
                    >
                      {item.label}
                    </a>
                  ))}
                </div>
              </div>
            )}
          </div>

          {NAV_LINKS.map(({ label, href }) => (
            <a
              key={label}
              className={cn(linkClass, "shrink-0")}
              href={href}
              onClick={(e) => e.preventDefault()}
            >
              {label}
            </a>
          ))}

          <a
            className={cn(linkClass, "inline-flex shrink-0 items-center gap-1")}
            href={DOCS_HREF}
            rel="noopener noreferrer"
            target="_blank"
          >
            Dokumentasi
            <ArrowUpRightIcon aria-hidden className="size-3.5 shrink-0" />
            <span className="sr-only"> (terbuka di tab baru)</span>
          </a>
        </nav>

        <div className="flex min-w-0 flex-1 items-center justify-end gap-1 sm:gap-2">
          <div className="hidden items-center gap-1 sm:gap-2 md:flex">
            <Button asChild variant="outline">
              <a href={loginHref} className="font-mono uppercase">
                Masuk
              </a>
            </Button>
            <Button asChild variant="default">
              <a href={loginHref} className="font-mono uppercase">
                Mulai Gratis
              </a>
            </Button>
          </div>
          <button
            aria-expanded={menuOpen}
            aria-label={menuOpen ? "Tutup menu" : "Buka menu"}
            className="ms-0.5 shrink-0 text-muted-foreground transition-colors hover:text-foreground md:hidden"
            type="button"
            onClick={() => setMenuOpen((v) => !v)}
          >
            <HamburgerIcon open={menuOpen} />
          </button>
        </div>
      </div>
    </header>
  )
}
