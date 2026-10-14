<?php

namespace Tests\Unit;

use App\Support\ChatPrompts;
use App\Support\Effort;
use PHPUnit\Framework\TestCase;

class EffortTest extends TestCase
{
    public function test_there_are_five_levels_in_ascending_order(): void
    {
        $this->assertSame(['low', 'medium', 'high', 'xhigh', 'ultra'], Effort::LEVELS);
    }

    public function test_unknown_or_missing_levels_fall_back_to_medium(): void
    {
        $this->assertSame('medium', Effort::normalize(null));
        $this->assertSame('medium', Effort::normalize('mega'));
        $this->assertSame('ultra', Effort::normalize('ultra'));
    }

    public function test_ultra_maps_to_the_highest_openrouter_level(): void
    {
        $this->assertSame('max', Effort::openRouter('ultra'));
        $this->assertSame('xhigh', Effort::openRouter('xhigh'));
        $this->assertSame('medium', Effort::openRouter('nonsense'));
    }

    public function test_every_level_other_than_medium_adds_a_distinct_directive(): void
    {
        $directives = array_map(Effort::directive(...), ['low', 'high', 'xhigh', 'ultra']);

        $this->assertSame('', Effort::directive('medium'));
        $this->assertCount(4, array_unique($directives));
        $this->assertNotContains('', $directives);
    }

    public function test_system_prompt_includes_mode_and_effort(): void
    {
        $prompt = ChatPrompts::system('coding', 'xhigh');

        $this->assertStringContainsString('Mode coding', $prompt);
        $this->assertStringContainsString('effort sangat tinggi', $prompt);
        $this->assertSame(ChatPrompts::system('chat'), ChatPrompts::system('chat', 'medium'));
    }
}
