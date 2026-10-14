import type { ComponentType, SVGProps } from "react"

import { VercelWordmark } from "@/components/ui/svgs/vercelWordmark"
import { VercelWordmarkDark } from "@/components/ui/svgs/vercelWordmarkDark"
import { StripeWordmark } from "@/components/ui/svgs/stripeWordmark"
import { StripeWordmarkDark } from "@/components/ui/svgs/stripeWordmarkDark"
import { SupabaseWordmarkLight } from "@/components/ui/svgs/supabaseWordmarkLight"
import { SupabaseWordmarkDark } from "@/components/ui/svgs/supabaseWordmarkDark"
import { GithubWordmarkLight } from "@/components/ui/svgs/githubWordmarkLight"
import { GithubWordmarkDark } from "@/components/ui/svgs/githubWordmarkDark"
import { CursorWordmarkLight } from "@/components/ui/svgs/cursorWordmarkLight"
import { CursorWordmarkDark } from "@/components/ui/svgs/cursorWordmarkDark"
import { NeonWordmarkLight } from "@/components/ui/svgs/neonWordmarkLight"
import { NeonWordmarkDark } from "@/components/ui/svgs/neonWordmarkDark"
import { ConvexWordmarkLight } from "@/components/ui/svgs/convexWordmarkLight"
import { ConvexWordmarkDark } from "@/components/ui/svgs/convexWordmarkDark"
import { clerkWordmarkLight } from "@/components/ui/svgs/clerkWordmarkLight"
import { clerkWordmarkDark } from "@/components/ui/svgs/clerkWordmarkDark"
import { SanityWordmarkLight } from "@/components/ui/svgs/sanityWordmarkLight"
import { SanityWordmarkDark } from "@/components/ui/svgs/sanityWordmarkDark"
import { MongodbWordmarkLight } from "@/components/ui/svgs/mongodbWordmarkLight"
import { MongodbWordmarkDark } from "@/components/ui/svgs/mongodbWordmarkDark"
import { TursoWordmarkLight } from "@/components/ui/svgs/tursoWordmarkLight"
import { TursoWordmarkDark } from "@/components/ui/svgs/tursoWordmarkDark"
import { AstroWordmarkLight } from "@/components/ui/svgs/astroWordmarkLight"
import { AstroWordmarkDark } from "@/components/ui/svgs/astroWordmarkDark"

export type SvgComponent = ComponentType<SVGProps<SVGSVGElement>>

export type Partner = {
  name: string
  light: SvgComponent
  dark?: SvgComponent
}

/**
 * Cell `i` rotates through partners `i`, `i + CELL_COUNT`, ...
 * so the first row of six is what renders on first paint.
 */
export const partners: Partner[] = [
  { name: "Vercel",   light: VercelWordmark,        dark: VercelWordmarkDark },
  { name: "Supabase", light: SupabaseWordmarkLight, dark: SupabaseWordmarkDark },
  { name: "Stripe",   light: StripeWordmark,        dark: StripeWordmarkDark },
  { name: "GitHub",   light: GithubWordmarkLight,   dark: GithubWordmarkDark },
  { name: "Clerk",    light: clerkWordmarkLight,    dark: clerkWordmarkDark },
  { name: "Cursor",   light: CursorWordmarkLight,   dark: CursorWordmarkDark },
  { name: "Astro",    light: AstroWordmarkLight,    dark: AstroWordmarkDark },
  { name: "Neon",     light: NeonWordmarkLight,     dark: NeonWordmarkDark },
  { name: "Convex",   light: ConvexWordmarkLight,   dark: ConvexWordmarkDark },
  { name: "MongoDB",  light: MongodbWordmarkLight,  dark: MongodbWordmarkDark },
  { name: "Sanity",   light: SanityWordmarkLight,   dark: SanityWordmarkDark },
  { name: "Turso",    light: TursoWordmarkLight,    dark: TursoWordmarkDark },
]
