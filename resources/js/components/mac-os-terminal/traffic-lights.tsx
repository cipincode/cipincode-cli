export default function TrafficLights() {
  return (
    <div className="absolute left-3 flex gap-2" aria-hidden>
      <button
        type="button"
        tabIndex={-1}
        className="size-3 shrink-0 rounded-full bg-[#ff5f57] ring-1 ring-black/10"
      />
      <button
        type="button"
        tabIndex={-1}
        className="size-3 shrink-0 rounded-full bg-[#febc2e] ring-1 ring-black/10"
      />
      <button
        type="button"
        tabIndex={-1}
        className="size-3 shrink-0 rounded-full bg-[#28c840] ring-1 ring-black/10"
      />
    </div>
  )
}
