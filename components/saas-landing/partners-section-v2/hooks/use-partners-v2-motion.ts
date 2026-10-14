"use client"

import { type RefObject, useLayoutEffect } from "react"

import gsap from "gsap"
import { ScrollTrigger } from "gsap/ScrollTrigger"

import { motion } from "../lib/config"

gsap.registerPlugin(ScrollTrigger)

type Cell = {
  el: HTMLElement
  layers: HTMLElement[]
  current: number
  busy: boolean
}

export function usePartnersV2Motion(rootRef: RefObject<HTMLElement | null>) {
  useLayoutEffect(() => {
    const root = rootRef.current
    if (!root) return

    const mm = gsap.matchMedia()

    mm.add("(prefers-reduced-motion: no-preference)", () => {
      const q = (selector: string) =>
        Array.from(root.querySelectorAll<HTMLElement>(selector))

      const words = q("[data-pv2-word]")
      const description = q("[data-pv2-description]")
      const linesX = q('[data-pv2-line="x"]')
      const linesY = q('[data-pv2-line="y"]')
      const marks = q("[data-pv2-mark]")
      const cells: Cell[] = q("[data-pv2-cell]").map((el) => ({
        el,
        layers: Array.from(el.querySelectorAll<HTMLElement>("[data-pv2-layer]")),
        current: 0,
        busy: false,
      }))
      const allLayers = cells.flatMap((cell) => cell.layers)
      const firstLogos = cells
        .map((cell) => cell.layers[0]?.querySelector<HTMLElement>("[data-pv2-logo]"))
        .filter((el): el is HTMLElement => Boolean(el))

      gsap.set(words, { yPercent: 110 })
      gsap.set(description, { autoAlpha: 0, y: 16 })
      gsap.set(linesX, { scaleX: 0, transformOrigin: "left center" })
      gsap.set(linesY, { scaleY: 0, transformOrigin: "center top" })
      gsap.set(marks, { autoAlpha: 0, scale: 0.3, rotate: -90 })
      gsap.set(firstLogos, { autoAlpha: 0, yPercent: 40, filter: "blur(8px)" })
      root.dataset.pv2Ready = ""

      // Touch browsers keep :hover on the last tapped element, which would
      // freeze that cell, so only honor hover on devices that can hover.
      const canHover = window.matchMedia("(hover: hover)").matches

      let inView = false
      let introDone = false
      let lastCell = -1
      let nextSwap: gsap.core.Tween | null = null

      // Swap one random cell to its next logo. Skips the cell that just
      // swapped, any cell mid-swap, and the cell under the pointer.
      const swap = () => {
        const candidates = cells.filter(
          (cell, i) =>
            i !== lastCell &&
            !cell.busy &&
            cell.layers.length > 1 &&
            !(canHover && cell.el.matches(":hover"))
        )
        if (!candidates.length) return

        const cell = candidates[Math.floor(Math.random() * candidates.length)]
        lastCell = cells.indexOf(cell)

        const outgoing = cell.layers[cell.current]
        cell.current = (cell.current + 1) % cell.layers.length
        const incoming = cell.layers[cell.current]
        cell.busy = true

        const d = motion.swapDuration
        const tl = gsap.timeline({
          onComplete: () => {
            cell.busy = false
          },
        })

        tl.to(
          outgoing,
          {
            yPercent: -60,
            autoAlpha: 0,
            filter: "blur(8px)",
            duration: d * 0.55,
            ease: "power3.in",
          },
          0
        ).fromTo(
          incoming,
          { yPercent: 60, autoAlpha: 0, filter: "blur(8px)" },
          {
            yPercent: 0,
            autoAlpha: 1,
            filter: "blur(0px)",
            duration: d * 0.75,
            ease: "power3.out",
          },
          d * 0.4
        )
      }

      const loop = () => {
        swap()
        nextSwap = gsap.delayedCall(motion.swapInterval, loop)
      }

      const start = () => {
        if (nextSwap || !introDone || !inView) return
        nextSwap = gsap.delayedCall(motion.swapInterval, loop)
      }

      const stop = () => {
        nextSwap?.kill()
        nextSwap = null
      }

      const intro = gsap.timeline({
        paused: true,
        defaults: { ease: "power3.out" },
        onComplete: () => {
          introDone = true
          start()
        },
      })

      // Title, then description, then the grid once both have mostly landed.
      intro
        .to(words, { yPercent: 0, duration: 0.9, stagger: 0.05, ease: "power4.out" })
        .to(description, { autoAlpha: 1, y: 0, duration: 0.7 }, 0.45)
        .addLabel("grid", "-=0.2")
        .to(linesX, { scaleX: 1, duration: 1.1, stagger: 0.07, ease: "expo.inOut" }, "grid")
        .to(linesY, { scaleY: 1, duration: 1.1, stagger: 0.07, ease: "expo.inOut" }, "grid+=0.1")
        .to(
          firstLogos,
          {
            autoAlpha: 1,
            yPercent: 0,
            filter: "blur(0px)",
            duration: 0.8,
            stagger: { each: 0.07, from: "random" },
            clearProps: "transform,filter",
          },
          "grid+=0.5"
        )
        .to(
          marks,
          { autoAlpha: 1, scale: 1, rotate: 0, duration: 0.6, stagger: 0.06, ease: "back.out(2)" },
          "grid+=0.6"
        )

      ScrollTrigger.create({
        trigger: root,
        start: "top 80%",
        end: "bottom top",
        onEnter: () => intro.play(),
        onEnterBack: () => intro.play(),
        onToggle: (self) => {
          inView = self.isActive
          if (inView) start()
          else stop()
        },
      })

      // Swap timelines are created after this callback returns, so
      // matchMedia's revert doesn't track them — clean them up by hand.
      return () => {
        stop()
        gsap.killTweensOf(allLayers)
        gsap.set(allLayers, { clearProps: "all" })
        delete root.dataset.pv2Ready
      }
    })

    return () => mm.revert()
  }, [rootRef])
}
